package main

import (
	"image/color"
	"testing"
)

func TestParseValidHexColor(t *testing.T) {
	colors := []string{
		"#ffffff",
		"#0D1117",
		"#eb4034",
		"#00ACD7",
	}

	for _, c := range colors {
		if _, err := ParseHexColor(c); err != nil {
			t.Errorf("Parse Error: %s", c)
		}
	}
}

func TestParseInvalidInvalidHexColor(t *testing.T) {
	colors := []string{
		"invalid c test",
		"ffffff",
		"abababab",
		"#12345aa",
		"#00",
	}

	for _, c := range colors {
		if _, err := ParseHexColor(c); err == nil {
			t.Errorf("Parse Error: %s", c)
		}
	}
}

func TestIsNoBackgroundColor(t *testing.T) {
	noBgColor := color.RGBA{}
	if !IsNoBackgroundColor(noBgColor) {
		t.Fatalf("Actual false, Expected: true")
	}

	bgColor := color.RGBA{R: 255, G: 0, B: 0, A: 0}
	if IsNoBackgroundColor(bgColor) {
		t.Fatalf("Actual true, Expected: false")
	}
}
