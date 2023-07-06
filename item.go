package flexpdf

import (
	"image/color"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

type flexItemExtender interface {
	FlexItem
	drawContent(*gopdf.GoPdf, rect, int) error
	getContentSize(*gopdf.GoPdf) (*size, error)
}

type flexItemCommon[T flexItemExtender] struct {
	self            T
	Width           float64
	Height          float64
	FlexGrow        float64
	FlexShrink      float64
	BackgroundColor color.Color
	Border          Border
	Margin          Spacing
	Padding         Spacing
}

func (c *flexItemCommon[T]) init(self T) {
	c.self = self
	c.Width = -1
	c.Height = -1
	c.FlexGrow = 0
	c.FlexShrink = 1
	c.BackgroundColor = nil
	c.Border = UniformedBorder(nil, BorderStyleSolid, 0) // TODO None
}

func (c *flexItemCommon[T]) SetWidth(w float64) T {
	c.Width = w
	return c.self
}
func (c *flexItemCommon[T]) SetHeight(h float64) T {
	c.Height = h
	return c.self
}
func (c *flexItemCommon[T]) SetSize(w, h float64) T {
	c.Width = w
	c.Height = h
	return c.self
}
func (b *flexItemCommon[T]) SetBackgroundColor(c color.Color) T {
	b.BackgroundColor = c
	return b.self
}
func (c *flexItemCommon[T]) SetBorder(border Border) T {
	c.Border = border
	return c.self
}
func (c *flexItemCommon[T]) SetMargin(margin Spacing) T {
	c.Margin = margin
	return c.self
}
func (c *flexItemCommon[T]) SetPadding(padding Spacing) T {
	c.Padding = padding
	return c.self
}

func (c *flexItemCommon[T]) draw(pdf *gopdf.GoPdf, marginBox rect, depth int) error {
	borderBox := marginBox
	borderBox.x += c.Margin.Left
	borderBox.w -= c.Margin.Left + c.Margin.Right
	borderBox.y += c.Margin.Top
	borderBox.h -= c.Margin.Top + c.Margin.Bottom
	paddingBox := borderBox
	paddingBox.x += c.Border.Left.Width
	paddingBox.w -= c.Border.Left.Width + c.Border.Right.Width
	paddingBox.y += c.Border.Top.Width
	paddingBox.h -= c.Border.Top.Width + c.Border.Bottom.Width
	contentBox := paddingBox
	contentBox.x += c.Padding.Left
	contentBox.w -= c.Padding.Left + c.Padding.Right
	contentBox.y += c.Padding.Top
	contentBox.h -= c.Padding.Top + c.Padding.Bottom

	// 背景色
	if c.BackgroundColor != nil && borderBox.w >= 0 && borderBox.h >= 0 {
		if err := setColor(pdf, c.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(borderBox.x, borderBox.y, borderBox.x+borderBox.w, borderBox.y+borderBox.h, "F", 0, 0); err != nil {
			return errors.Wrap(err, "rectangle")
		}
	}
	if err := c.self.drawContent(pdf, contentBox, depth); err != nil {
		return err
	}
	if err := c.Border.draw(pdf, borderBox); err != nil {
		return err
	}
	return nil
}
func (c *flexItemCommon[T]) getPreferredSize(pdf *gopdf.GoPdf) (*size, error) {
	ps, err := c.self.getContentSize(pdf)
	if err != nil {
		return nil, err
	}
	if c.Width >= 0 {
		ps.w = c.Width
	}
	if c.Height >= 0 {
		ps.h = c.Height
	}
	ps.w += c.Margin.Right + c.Margin.Left
	ps.h += c.Margin.Top + c.Margin.Bottom

	return ps, nil
}
