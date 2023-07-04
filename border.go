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
	Top    BorderPart
	Right  BorderPart
	Bottom BorderPart
	Left   BorderPart
}

func UniformedBorder(co color.Color, style BorderStyle, width float64) Border {
	p := BorderPart{co, width, style}
	return Border{p, p, p, p}
}

func (b *Border) draw(pdf *gopdf.GoPdf, r rect) error {
	if err := b.Top.draw(pdf, r.x, r.y, r.x+r.w, r.y); err != nil {
		return err
	}
	if err := b.Right.draw(pdf, r.x+r.w, r.y, r.x+r.w, r.y+r.h); err != nil {
		return err
	}
	if err := b.Bottom.draw(pdf, r.x, r.y+r.h, r.x+r.w, r.y+r.h); err != nil {
		return err
	}
	if err := b.Left.draw(pdf, r.x, r.y, r.x, r.y+r.h); err != nil {
		return err
	}
	return nil
}

type BorderPart struct {
	Color color.Color
	Width float64
	Style BorderStyle
}

func (p *BorderPart) draw(pdf *gopdf.GoPdf, x1, y1, x2, y2 float64) error {
	if p.Color != nil {
		if err := setColor(pdf, p.Color); err != nil {
			return err
		}
	}
	switch p.Style {
	case BorderStyleDashed:
		pdf.SetLineType("dashed")
	case BorderStyleDotted:
		pdf.SetLineType("dotted")
	case BorderStyleSolid:
		pdf.SetLineType("")
	default:
		panic(p.Style)
	}

	pdf.SetLineWidth(p.Width)
	pdf.Line(x1, y1, x2, y2)
	return nil
}
