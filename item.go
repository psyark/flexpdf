package flexpdf

import "image/color"

type flexItemCommon[T any] struct {
	self            *T
	Width           float64
	Height          float64
	FlexGrow        float64
	FlexShrink      float64
	BackgroundColor color.Color
	Border          Border
	Margin          Spacing
	Padding         Spacing
}

func (c *flexItemCommon[T]) init(self *T) {
	c.self = self
	c.Width = -1
	c.Height = -1
	c.FlexGrow = 0
	c.FlexShrink = 1
	c.BackgroundColor = nil
	c.Border = UniformedBorder(nil, BorderStyleSolid, 0) // TODO None
}

func (c *flexItemCommon[T]) SetWidth(w float64) *T {
	c.Width = w
	return c.self
}
func (c *flexItemCommon[T]) SetHeight(h float64) *T {
	c.Height = h
	return c.self
}
func (c *flexItemCommon[T]) SetSize(w, h float64) *T {
	c.Width = w
	c.Height = h
	return c.self
}
func (b *flexItemCommon[T]) SetBackgroundColor(c color.Color) *T {
	b.BackgroundColor = c
	return b.self
}
func (c *flexItemCommon[T]) SetBorder(border Border) *T {
	c.Border = border
	return c.self
}
