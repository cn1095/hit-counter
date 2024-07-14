package badge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/math/fixed"
)

func TestFontM_FixedPointToPoint(t *testing.T) {
	tests := []struct {
		name   string
		drawer *fontM
		input  fixed.Int26_6
		output float64
	}{
		{name: "0", drawer: &fontM{}, input: fixed.I(0), output: 0},
		{name: "100", drawer: &fontM{}, input: fixed.I(100), output: 100},
		{name: "-255", drawer: &fontM{}, input: fixed.I(-255), output: -255},
		{name: "20.25", drawer: &fontM{}, input: fixed.I(20), output: 20},
		{name: "-100.875", drawer: &fontM{}, input: fixed.I(-100), output: -100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.drawer.fixedToPoint(tc.input))
			fmt.Printf("%s = %f\n", tc.input.String(), tc.output)
		})
	}
}

func TestFontM_MeasureString(t *testing.T) {
	tests := []struct {
		name   string
		drawer fontDrawer
		input  string
		output float64
	}{
		{name: "vera_sans_1", drawer: veraSansDrawer, input: "allan hi", output: 53},
		{name: "verdana_1", drawer: verdanaDrawer, input: "allan hi", output: 51},
		{name: "vera_sans_2", drawer: veraSansDrawer, input: "hits", output: 34},
		{name: "verdana_2", drawer: verdanaDrawer, input: "hits", output: 30},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.drawer.measureString(tc.input))
			fmt.Printf("%s = %f \n", tc.input, tc.drawer.measureString(tc.input))
		})
	}

}

func TestFontM_FontSize(t *testing.T) {
	tests := []struct {
		name   string
		drawer fontDrawer
		output int
	}{
		{name: "ok", drawer: &fontM{fontSize: fontSize}, output: 11},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.drawer.getFontSize())
		})
	}
}

func TestFontM_FontFamily(t *testing.T) {
	tests := []struct {
		name   string
		drawer fontDrawer
		output string
	}{
		{name: "ok", drawer: &fontM{fontFamily: fontFamilyVeraSans}, output: fontFamilyVeraSans},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.drawer.getFontFamily())
		})
	}
}

func TestGetFontDrawer(t *testing.T) {
	tests := []struct {
		name   string
		input  FontType
		output fontDrawer
		isErr  bool
	}{
		{name: "fail", input: 0, isErr: true},
		{name: "vera sans", input: VeraSans, isErr: false, output: veraSansDrawer},
		{name: "verdana", input: Verdana, isErr: false, output: verdanaDrawer},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			drawer, err := getFontDrawer(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.Equal(t, tc.output, drawer)
			}
		})
	}
}
