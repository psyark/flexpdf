package flexpdf

import (
	"image/color"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

// Text はテキストを扱うエレメントです
// TODO family, size, color, text は Spanのスライスにする
type Text struct {
	flexItemCommon

	FontFamily string
	FontSize   float64
	Text       string
	Color      color.Color
	LineHeight float64
}

func NewText(family string, size float64, text string) *Text {
	t := &Text{}
	t.Width = -1
	t.Height = -1

	t.FontFamily = family
	t.FontSize = size
	t.Text = text
	return t
}
func (t *Text) SetWidth(w float64) *Text {
	t.Width = w
	return t
}
func (t *Text) SetHeight(h float64) *Text {
	t.Height = h
	return t
}
func (t *Text) SetSize(w, h float64) *Text {
	return t.SetWidth(w).SetHeight(h)
}
func (t *Text) SetBackgroundColor(c color.Color) *Text {
	t.BackgroundColor = c
	return t
}

func (t *Text) draw(pdf *gopdf.GoPdf, r rect, depth int) error {
	log.Printf("%sText.draw(r=%v, t=%q)\n", strings.Repeat("  ", depth), r, t.Text)

	// 背景色
	if t.BackgroundColor != nil && r.w != 0 && r.h != 0 {
		if err := setColor(pdf, t.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(r.x, r.y, r.x+r.w, r.y+r.h, "F", 0, 0); err != nil {
			return errors.Wrap(err, "rectangle")
		}
	}

	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return errors.Wrap(err, "setFont")
	}

	{
		c := t.Color
		if c == nil {
			c = color.Black
		}
		if err := setColor(pdf, c); err != nil {
			return errors.Wrap(err, "setColor")
		}
	}

	pdf.SetXY(r.x, r.y)
	// TODO 幅が小さすぎる場合に無限ループになるのを抑制
	if r.w < 20 {
		return nil
	}

	lines, err := pdf.SplitTextWithWordWrap(t.Text, r.w)
	if err != nil {
		return errors.Wrap(err, "splitTextWithWordWrap")
	}

	for _, line := range lines {
		if err := pdf.MultiCell(&gopdf.Rect{W: r.w, H: r.h}, line); err != nil {
			return errors.Wrap(err, "multiCell")
		}
		if t.LineHeight != 0 && t.LineHeight != 1 {
			pdf.Br((t.LineHeight - 1) * t.FontSize)
			pdf.SetX(r.x)
		}
	}

	if err := t.Border.draw(pdf, r); err != nil {
		return err
	}

	return nil
}
func (t *Text) getPreferredSize(pdf *gopdf.GoPdf) (*size, error) {
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return nil, err
	}

	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return nil, err
	}

	ps := &size{
		h: t.FontSize * float64(len(lines)),
	}
	if t.LineHeight != 0 {
		ps.h += (t.FontSize * (t.LineHeight - 1)) * float64(len(lines)-1)
	}

	for _, line := range lines {
		w, err := pdf.MeasureTextWidth(line)
		if err != nil {
			return nil, err
		}
		if ps.w < w {
			ps.w = w
		}
	}

	if t.Width >= 0 {
		ps.w = t.Width
	}
	if t.Height >= 0 {
		ps.h = t.Height
	}

	return ps, nil
}

type TextSpan struct {
}
