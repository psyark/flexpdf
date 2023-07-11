package flexpdf

import (
	"image/color"
	"log"

	"github.com/signintech/gopdf"
)

var (
	_ flexItemContent = &Box{}
	_ flexItemContent = &Text{}
)

// flexItemContent は
type flexItemContent interface {
	FlexItem
	drawContent(*gopdf.GoPdf, rect) error
	getContentSize(pdf *gopdf.GoPdf, contentBoxMax size) (size, error)
}

type flexItemCommon[T flexItemContent] struct {
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

func (c *flexItemCommon[T]) getFlexGrow() float64 {
	return c.FlexGrow
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
func (c *flexItemCommon[T]) SetFlexGrow(g float64) T {
	c.FlexGrow = g
	return c.self
}
func (c *flexItemCommon[T]) SetFlexShrink(s float64) T {
	c.FlexShrink = s
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
func (*flexItemCommon[T]) parseSpacing(values ...float64) Spacing {
	switch len(values) {
	case 0:
		return Spacing{}
	case 1: // TRBL
		return Spacing{
			Top:    values[0],
			Right:  values[0],
			Bottom: values[0],
			Left:   values[0],
		}
	case 2: // TB | RL
		return Spacing{
			Top:    values[0],
			Bottom: values[0],
			Right:  values[1],
			Left:   values[1],
		}
	case 3: // T | RL | B
		return Spacing{
			Top:    values[0],
			Right:  values[1],
			Left:   values[1],
			Bottom: values[2],
		}
	default: // T | R | B | L
		return Spacing{
			Top:    values[0],
			Right:  values[1],
			Bottom: values[2],
			Left:   values[3],
		}
	}
}
func (c *flexItemCommon[T]) SetMargin(values ...float64) T {
	c.Margin = c.parseSpacing(values...)
	return c.self
}
func (c *flexItemCommon[T]) SetPadding(values ...float64) T {
	c.Padding = c.parseSpacing(values...)
	return c.self
}

func (c *flexItemCommon[T]) draw(pdf *gopdf.GoPdf, marginBox rect) (err error) {
	defer wrap(&err, "common.draw")

	borderBox := marginBox.shrink(c.Margin)
	paddingBox := borderBox.shrink(c.Border.Width)
	contentBox := paddingBox.shrink(c.Padding)

	// 背景色
	if c.BackgroundColor != nil && borderBox.w > 0 && borderBox.h > 0 {
		if err := setColor(pdf, c.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(borderBox.x, borderBox.y, borderBox.x+borderBox.w, borderBox.y+borderBox.h, "F", 0, 0); err != nil {
			return err
		}
	}
	if err := c.self.drawContent(pdf, contentBox); err != nil {
		return err
	}
	if err := c.Border.draw(pdf, borderBox); err != nil {
		return err
	}
	return nil
}
func (c *flexItemCommon[T]) getPreferredSize(pdf *gopdf.GoPdf, marginBoxMax size) (size, error) {
	contentBoxMax := marginBoxMax.shrink(c.Margin).shrink(c.Border.Width).shrink(c.Padding)

	ps, err := c.self.getContentSize(pdf, contentBoxMax)
	if err != nil {
		return size{}, err
	}

	// TODO もしWidthが指定されてるなら、そこからMargin, Border, Paddingを引いてから
	// getContentSizeにwidthを渡す
	if c.Width >= 0 {
		ps.w = c.Width
	}
	if c.Height >= 0 {
		ps.h = c.Height
	}

	for _, space := range []Spacing{c.Margin, c.Border.Width, c.Padding} {
		ps = ps.expand(space)
	}

	if ps.h > 600 {
		log.Println("🍣", ps)
	}

	return ps, nil
}
