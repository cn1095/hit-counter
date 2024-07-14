package badge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	tests := []struct {
		name  string
		isErr bool
	}{
		{
			name:  "success",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewWriter()
			assert.Equal(t, tc.isErr, err != nil)
		})
	}
}

// svg parser https://www.rapidtables.com/web/tools/svg-viewer-editor.html
func TestBadgeWriter_RenderFlatBadge(t *testing.T) {
	tests := []struct {
		name   string
		input  Badge
		outupt []byte
		isErr  bool
	}{
		{
			name: "verasans-flat",
			input: Badge{
				FontType:             VeraSans,
				LeftText:             "verasans-flat",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "10 / 23234",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "3",
				YRadius:              "3",
			},
		},
		{
			name: "verasans-round",
			input: Badge{
				FontType:             VeraSans,
				LeftText:             "verasans-round",
				LeftTextColor:        "#1E9268",
				LeftBackgroundColor:  "#252050",
				RightText:            "1 / 202",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "auto",
				YRadius:              "auto",
			},
		},
		{
			name: "verdana-flat",
			input: Badge{
				FontType:             Verdana,
				LeftText:             "verdana-flat",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "10 / 23234",
				RightTextColor:       "#E25C9F",
				RightBackgroundColor: "#502038",
				XRadius:              "3",
				YRadius:              "3",
			},
		},
		{
			name: "verdana-round",
			input: Badge{
				FontType:             Verdana,
				LeftText:             "verdand-round",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "1 / 202",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "auto",
				YRadius:              "auto",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			writer, err := NewWriter()
			assert.NoError(t, err)
			result, err := writer.RenderFlatBadge(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			fmt.Println(string(result))
		})
	}
}

// svg parser https://www.rapidtables.com/web/tools/svg-viewer-editor.html
func TestBadgeWriter_RenderIconBadge(t *testing.T) {
	tests := []struct {
		name      string
		input     Badge
		iconName  string
		iconColor string
		outupt    []byte
		isErr     bool
	}{
		{
			name: "verasans-flat",
			input: Badge{
				FontType:             VeraSans,
				LeftText:             "verasans-flat",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "10 / 23234",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "3",
				YRadius:              "3",
			},
			iconName:  "appveyor.svg",
			iconColor: "#00B3E0",
		},
		{
			name: "verasans-round",
			input: Badge{
				FontType:             VeraSans,
				LeftText:             "verasans-round",
				LeftTextColor:        "#1E9268",
				LeftBackgroundColor:  "#252050",
				RightText:            "1 / 202",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "auto",
				YRadius:              "auto",
			},
			iconName:  "appveyor.svg",
			iconColor: "#3c5688",
		},
		{
			name: "verdana-flat",
			input: Badge{
				FontType:             Verdana,
				LeftText:             "verdana-flat",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "10 / 23234",
				RightTextColor:       "#E25C9F",
				RightBackgroundColor: "#502038",
				XRadius:              "3",
				YRadius:              "3",
			},
			iconName:  "amazon.svg",
			iconColor: "#0000ff",
		},
		{
			name: "verdana-round",
			input: Badge{
				FontType:             Verdana,
				LeftText:             "verdand-round",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "1 / 202",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "auto",
				YRadius:              "auto",
			},
			iconName:  "babel.svg",
			iconColor: "#109556",
		},
		{
			name: "verdana-hits",
			input: Badge{
				FontType:             Verdana,
				LeftText:             "hits",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            "1 / 202",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#4c1",
				XRadius:              "auto",
				YRadius:              "auto",
			},
			iconName:  "aircall.svg",
			iconColor: "#ffffff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			writer, err := NewWriter()
			assert.NoError(t, err)
			result, err := writer.RenderIconBadge(tc.input, tc.iconName, tc.iconColor)
			assert.Equal(t, tc.isErr, err != nil)
			fmt.Println(string(result))
		})
	}
}

func TestGetIconsMap(t *testing.T) {
	assert.Equal(t, len(iconsMap), len(GetIconsMap()))
	for k, _ := range GetIconsMap() {
		fmt.Println(k)
	}
}
