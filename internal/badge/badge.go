package badge

import (
	"github.com/gjbae1212/hit-counter/pkg/badge"
)

// GenerateBadge makes Flat-Badge struct which is used go-counter-badge/badge.
func GenerateBadge(leftText, leftBgColor, rightText, rightBgColor string, edgeFlat bool) badge.Badge {
	flatBadge := badge.Badge{
		FontType:             badge.Verdana,
		LeftText:             leftText,
		LeftTextColor:        "#fff",
		LeftBackgroundColor:  leftBgColor,
		RightText:            rightText,
		RightTextColor:       "#fff",
		RightBackgroundColor: rightBgColor,
	}
	if edgeFlat {
		flatBadge.XRadius = "0"
		flatBadge.YRadius = "0"
	} else {
		flatBadge.XRadius = "3"
		flatBadge.YRadius = "3"
	}
	return flatBadge
}
