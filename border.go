package flexpdf

import (
	"image/color"

	"github.com/signintech/gopdf"
)

type BorderStyle int

const (
	BorderStyleSolid BorderStyle = iota
	BorderStyleDashed
	BorderStyleDotted
)

type Border struct {
	Color TRBL[color.Color]
	Width Spacing
	Style TRBL[BorderStyle]
}

func UniformedBorder(col color.Color, style BorderStyle, width float64) Border {
	return Border{
		Color: TRBL[color.Color]{col, col, col, col},
		Width: Spacing{width, width, width, width},
		Style: TRBL[BorderStyle]{style, style, style, style},
	}
}

func (b *Border) draw(pdf *gopdf.GoPdf, r rect) (err error) {
	defer wrap(&err, "border.draw")

	r.x += b.Width.Left / 2
	r.y += b.Width.Top / 2
	r.w -= (b.Width.Left + b.Width.Right) / 2
	r.h -= (b.Width.Top + b.Width.Bottom) / 2

	if err := b.drawPart(pdf, r.x, r.y, r.x+r.w, r.y, b.Color.Top, b.Width.Top, b.Style.Top); err != nil {
		return err
	}
	if err := b.drawPart(pdf, r.x+r.w, r.y, r.x+r.w, r.y+r.h, b.Color.Right, b.Width.Right, b.Style.Right); err != nil {
		return err
	}
	if err := b.drawPart(pdf, r.x, r.y+r.h, r.x+r.w, r.y+r.h, b.Color.Bottom, b.Width.Bottom, b.Style.Bottom); err != nil {
		return err
	}
	if err := b.drawPart(pdf, r.x, r.y, r.x, r.y+r.h, b.Color.Left, b.Width.Left, b.Style.Left); err != nil {
		return err
	}
	return nil
}

func (*Border) drawPart(pdf *gopdf.GoPdf, x1, y1, x2, y2 float64, col color.Color, width float64, style BorderStyle) error {
	if col != nil && width > 0 {
		if err := setColor(pdf, col); err != nil {
			return err
		}

		switch style {
		case BorderStyleDashed:
			pdf.SetLineType("dashed")
		case BorderStyleDotted:
			pdf.SetLineType("dotted")
		case BorderStyleSolid:
			pdf.SetLineType("")
		default:
			panic(style)
		}

		pdf.SetLineWidth(width)
		pdf.Line(x1, y1, x2, y2)
	}
	return nil
}
