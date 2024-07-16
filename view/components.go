package view

import (
	"fmt"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

// Components
func radioButton(label string, checked bool) string {
	if checked {
		return radioButtonStyle.Render("(o) " + label)
	}
	return "( ) " + label
}

// Utils

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
