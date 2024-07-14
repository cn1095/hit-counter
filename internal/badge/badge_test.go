package badge

import (
	"reflect"
	"testing"

	"github.com/gjbae1212/hit-counter/pkg/badge"
	"github.com/stretchr/testify/assert"
)

func TestGenerateBadge(t *testing.T) {
	tests := []struct {
		name         string
		leftText     string
		leftBgColor  string
		rightText    string
		rightBgColor string
		edgeFlat     bool
		output       badge.Badge
	}{
		{name: "not-edge",
			leftText:     "allan",
			leftBgColor:  "#555",
			rightText:    " 0 / 10 ",
			rightBgColor: "#79c83d",
			edgeFlat:     false,
			output: badge.Badge{
				FontType:             badge.Verdana,
				LeftText:             "allan",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            " 0 / 10 ",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#79c83d",
				XRadius:              "3",
				YRadius:              "3",
			},
		},
		{name: "edge",
			leftText:     "allan",
			leftBgColor:  "#555",
			rightText:    " 0 / 10 ",
			rightBgColor: "#79c83d",
			edgeFlat:     true,
			output: badge.Badge{
				FontType:             badge.Verdana,
				LeftText:             "allan",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            " 0 / 10 ",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#79c83d",
				XRadius:              "0",
				YRadius:              "0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bg := GenerateBadge(tc.leftText, tc.leftBgColor, tc.rightText, tc.rightBgColor, tc.edgeFlat)
			assert.True(t, reflect.DeepEqual(tc.output, bg))
		})
	}
}
