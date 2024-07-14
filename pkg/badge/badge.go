package badge

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"github.com/gjbae1212/hit-counter/pkg/badge/internal/assets"
	perrors "github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	defaultBadgeHeight = float64(20)
	defaultIconWidth   = float64(15)
	defaultIconHeight  = float64(15)
	defaultIconX       = float64(3)
	defaultIconY       = float64(2.5)
)

var (
	iconsMap = map[string]Icon{}
)

type Icon struct {
	Name   string
	Origin []byte
}

type Badge struct {
	FontType FontType

	LeftText            string
	LeftTextColor       string
	LeftBackgroundColor string

	RightText            string
	RightTextColor       string
	RightBackgroundColor string

	XRadius string
	YRadius string
}

// Writer is an interface generating Badge formatted SVG.
type Writer interface {
	RenderFlatBadge(b Badge) ([]byte, error)
	RenderIconBadge(b Badge, iconName, iconColor string) ([]byte, error)
}

type badgeWriter struct {
	tmplIconBadge *template.Template // icon-badge template
	tmplFlatBadge *template.Template // flat-badge template
}

// RenderFlatBadge renders Flat Badge formatted SVG to byte array.
func (fb *badgeWriter) RenderFlatBadge(b Badge) ([]byte, error) {
	drawer, err := getFontDrawer(b.FontType)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	// default dy
	dy := defaultBadgeHeight

	// set x,y radius
	fBadge := &flatBadge{FontFamily: drawer.getFontFamily(), FontSize: drawer.getFontSize()}
	fBadge.Rx = b.XRadius
	fBadge.Ry = b.YRadius

	// set left
	leftDx := drawer.measureString(b.LeftText)
	fBadge.Left = badge{
		Rect: rect{Color: color(b.LeftBackgroundColor), Bound: bound{
			Dx: leftDx,
			Dy: dy,
			X:  0,
			Y:  0,
		}},
		Text: text{Msg: b.LeftText, Color: color(b.LeftTextColor), Bound: bound{
			Dx: 0, // not use
			Dy: 0, // not use
			X:  leftDx/2.0 + 1,
			Y:  15,
		}},
	}

	// set right
	rightDx := drawer.measureString(b.RightText)
	fBadge.Right = badge{
		Rect: rect{Color: color(b.RightBackgroundColor), Bound: bound{
			Dx: rightDx,
			Dy: dy,
			X:  leftDx,
			Y:  0,
		}},
		Text: text{Msg: b.RightText, Color: color(b.RightTextColor), Bound: bound{
			Dx: 0, // not use
			Dy: 0, // not use
			X:  leftDx + rightDx/2.0 - 1,
			Y:  15,
		}},
	}

	// set dx, dy
	fBadge.Dy = defaultBadgeHeight
	fBadge.Dx = leftDx + rightDx

	buf := &bytes.Buffer{}
	if err := fb.tmplFlatBadge.Execute(buf, fBadge); err != nil {
		return nil, perrors.WithStack(err)
	}
	return buf.Bytes(), nil
}

// RenderIconBadge renders Icon Badge formatted SVG to byte array.
func (fb *badgeWriter) RenderIconBadge(b Badge, iconName, iconColor string) ([]byte, error) {
	if iconName == "" {
		return nil, perrors.WithStack(fmt.Errorf("[err] icon name empty"))
	}
	icon, ok := iconsMap[iconName]
	if !ok {
		return nil, perrors.WithStack(fmt.Errorf("[err] not found icon %s", iconName))
	}

	drawer, err := getFontDrawer(b.FontType)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	// fill icon color
	iconsvg := string(icon.Origin)
	if iconColor != "" {
		iconsvg = strings.Replace(iconsvg, "<svg", fmt.Sprintf("<svg fill=\"%s\" ", iconColor), 1)
	}

	// default dy
	dy := defaultBadgeHeight

	// set x,y radius
	iBadge := &iconBadge{FontFamily: drawer.getFontFamily(), FontSize: drawer.getFontSize()}
	iBadge.Rx = b.XRadius
	iBadge.Ry = b.YRadius

	// set icon
	iconDx := defaultIconWidth + 2*defaultIconX
	iBadge.Icon.Bound.X = defaultIconX
	iBadge.Icon.Bound.Y = defaultIconY
	iBadge.Icon.Bound.Dx = iconDx
	iBadge.Icon.Bound.Dy = defaultIconHeight
	iBadge.Icon.Base64 = base64.StdEncoding.EncodeToString([]byte(iconsvg))

	// set left
	leftDx := drawer.measureString(b.LeftText)
	iBadge.Left = badge{
		Rect: rect{Color: color(b.LeftBackgroundColor), Bound: bound{
			Dx: leftDx + iconDx,
			Dy: dy,
			X:  0, // not use
			Y:  0, // not use
		}},
		Text: text{Msg: b.LeftText, Color: color(b.LeftTextColor), Bound: bound{
			Dx: 0, // not use
			Dy: 0, // not use
			X:  leftDx/2.0 - 1 + iconDx,
			Y:  15,
		}},
	}

	// set right
	rightDx := drawer.measureString(b.RightText)
	iBadge.Right = badge{
		Rect: rect{Color: color(b.RightBackgroundColor), Bound: bound{
			Dx: rightDx,
			Dy: dy,
			X:  leftDx + iconDx,
			Y:  0,
		}},
		Text: text{Msg: b.RightText, Color: color(b.RightTextColor), Bound: bound{
			Dx: 0, // not use
			Dy: 0, // not use
			X:  iconDx + leftDx + rightDx/2.0,
			Y:  15,
		}},
	}

	// set dx, dy
	iBadge.Dy = defaultBadgeHeight
	iBadge.Dx = leftDx + rightDx + iconDx

	buf := &bytes.Buffer{}
	if err := fb.tmplIconBadge.Execute(buf, iBadge); err != nil {
		return nil, perrors.WithStack(err)
	}
	return buf.Bytes(), nil
}

// NewWriter returns Badge Writer.
func NewWriter() (Writer, error) {
	// make flat-badge template
	tmplFlatBadge, err := template.New("flat-badge").Parse(flatBadgeTemplate)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	tmplIconBadge, err := template.New("icon-badge").Parse(iconBadgeTemplate)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	writer := &badgeWriter{
		tmplFlatBadge: tmplFlatBadge,
		tmplIconBadge: tmplIconBadge,
	}
	return writer, nil
}

// GetIconsMap returns cloned iconsMap.
func GetIconsMap() map[string]Icon {
	cloned := make(map[string]Icon, len(iconsMap))
	for k, v := range iconsMap {
		cloned[k] = v
	}
	return cloned
}

func init() {
	entries, err := assets.Icons.ReadDir("icons")
	if err != nil {
		panic(perrors.WithStack(err))
	}

	iconNames := lo.Map(entries, func(item fs.DirEntry, index int) string {
		return item.Name()
	})

	for _, name := range iconNames {
		bin, err := assets.Icons.ReadFile("icons/" + name)
		if err != nil {
			panic(perrors.WithStack(err))
		}
		iconsMap[name] = Icon{Name: name, Origin: bin}
	}
}
