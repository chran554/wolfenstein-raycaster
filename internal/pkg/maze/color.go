package maze

import "image/color"

type Color struct {
	R, G, B float64
}

func NewColor(r, g, b float64) *Color {
	return &Color{R: r, G: g, B: b}
}

func NewColorFromByte(r, g, b byte) *Color {
	const byteNormalize = 1.0 / 255.0
	return &Color{R: float64(r) * byteNormalize, G: float64(g) * byteNormalize, B: float64(b) * byteNormalize}
}

func NewColorFromColor(color color.Color) *Color {
	const byteNormalize = 1.0 / 255.0

	r, g, b, _ := color.RGBA()
	return &Color{
		R: float64(r>>8) * byteNormalize,
		G: float64(g>>8) * byteNormalize,
		B: float64(b>>8) * byteNormalize,
	}
}

func (c *Color) RGBA() color.RGBA {
	r, g, b := c.Bytes()
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func (c *Color) FadeTo(to *Color, factor float64) *Color {
	return &Color{
		R: clamp(c.R+(to.R-c.R)*factor, 0.0, 1.0),
		G: clamp(c.G+(to.G-c.G)*factor, 0.0, 1.0),
		B: clamp(c.B+(to.B-c.B)*factor, 0.0, 1.0),
	}
}

func (c *Color) Mul(m *Color) *Color {
	return &Color{
		R: c.R * m.R,
		G: c.G * m.G,
		B: c.B * m.B,
	}
}

func (c *Color) Add(m *Color) *Color {
	return &Color{
		R: c.R + m.R,
		G: c.G + m.G,
		B: c.B + m.B,
	}
}

func (c *Color) Scale(t float64) *Color {
	return &Color{
		R: c.R * t,
		G: c.G * t,
		B: c.B * t,
	}
}

func (c *Color) Bytes() (r, g, b byte) {
	return byte(clamp(c.R, 0.0, 1.0) * 255.0),
		byte(clamp(c.G, 0.0, 1.0) * 255.0),
		byte(clamp(c.B, 0.0, 1.0) * 255.0)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
