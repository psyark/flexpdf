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
	flexItemCommon[*Text]

	FontFamily string
	FontSize   float64
	Text       string
	Color      color.Color
	LineHeight float64
}

func (t *Text) SetLineHeight(lineHeight float64) *Text {
	t.LineHeight = lineHeight
	return t
}

func NewText(family string, size float64, text string) *Text {
	t := &Text{
		FontFamily: family,
		FontSize:   size,
		Text:       text,
		LineHeight: 1,
	}
	t.flexItemCommon.init(t)
	return t
}

func (t *Text) drawContent(pdf *gopdf.GoPdf, r rect, depth int) error {
	log.Printf("%sText.draw(r=%v, t=%q)\n", strings.Repeat("  ", depth), r, t.Text)

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

	return nil
}
func (t *Text) getContentSize(pdf *gopdf.GoPdf) (*size, error) {
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return nil, err
	}

	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return nil, err
	}

	cs := &size{
		h: t.FontSize * (float64(len(lines)) + (t.LineHeight-1)*float64(len(lines)-1)),
	}

	for _, line := range lines {
		w, err := pdf.MeasureTextWidth(line)
		if err != nil {
			return nil, err
		}
		if cs.w < w {
			cs.w = w
		}
	}
	return cs, nil
}

type TextSpan struct {
}
