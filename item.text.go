package flexpdf

import (
	"image/color"
	"math"
	"strings"

	"github.com/signintech/gopdf"
)

type TextAlign int

const (
	TextAlignBegin TextAlign = iota
	TextAlignCenter
	TextAlignEnd
)

// Text はテキストを扱うエレメントです
// TODO family, size, color, text は Spanのスライスにする
type Text struct {
	flexItemCommon[*Text]

	Align TextAlign
	Runs  []*TextRun
}

type TextRun struct {
	Color      color.Color
	FontSize   float64
	FontFamily string
	LineHeight float64
	Text       string
}

func NewRun(text string) *TextRun {
	return &TextRun{
		Color:      color.Black,
		FontSize:   10,
		FontFamily: "",
		LineHeight: 1,
		Text:       text,
	}
}
func (r *TextRun) SetColor(c color.Color) *TextRun {
	r.Color = c
	return r
}
func (r *TextRun) SetFontSize(s float64) *TextRun {
	r.FontSize = s
	return r
}
func (r *TextRun) SetFontFamily(f string) *TextRun {
	r.FontFamily = f
	return r
}
func (r *TextRun) SetLineHeight(lineHeight float64) *TextRun {
	r.LineHeight = lineHeight
	return r
}
func (r *TextRun) SetText(t string) *TextRun {
	r.Text = t
	return r
}

// splitToNBR は改行コードのみを考慮して noBrRunのリストに分割します
func (r *TextRun) splitWithNewline() []noBrRun {
	nbrs := []noBrRun{}
	nbr := noBrRun{*r}
	for _, text := range strings.Split(r.Text, "\n") {
		nbr.Text = text
		nbrs = append(nbrs, nbr)
	}
	return nbrs
}

// 改行を含まない TextRun
type noBrRun struct {
	TextRun
}

func (r *noBrRun) size(pdf *gopdf.GoPdf) (size, error) {
	if err := pdf.SetFont(r.FontFamily, "", r.FontSize); err != nil {
		return size{}, err
	}
	w, err := pdf.MeasureTextWidth(r.Text)
	if err != nil {
		return size{}, err
	}
	return size{w: w, h: r.FontSize * r.LineHeight}, nil
}
func (r *noBrRun) splitWithWidth(pdf *gopdf.GoPdf, widthLimit float64) (*noBrRun, *noBrRun, error) {
	if widthLimit < 0 {
		return r, nil, nil
	}

	if err := pdf.SetFont(r.FontFamily, "", r.FontSize); err != nil {
		return nil, nil, err
	}

	runes := []rune(r.Text)
	for i := 1; i <= len(runes); i++ {
		w, err := pdf.MeasureTextWidth(string(runes[:i]))
		if err != nil {
			return nil, nil, err
		}
		if w > widthLimit {
			if i > 1 {
				i--
			}
			nbr1 := noBrRun{TextRun: r.TextRun}
			nbr1.Text = string(runes[:i])
			nbr2 := noBrRun{TextRun: r.TextRun}
			nbr2.Text = string(runes[i:])
			return &nbr1, &nbr2, nil
		}
	}
	return r, nil, nil
}

func (r *noBrRun) draw(pdf *gopdf.GoPdf) error {
	if err := pdf.SetFont(r.FontFamily, "", r.FontSize); err != nil {
		return err
	}
	if err := setColor(pdf, r.Color); err != nil {
		return err
	}

	s, err := r.size(pdf)
	if err != nil {
		return err
	}

	return pdf.Cell(&gopdf.Rect{W: s.w, H: s.h}, r.Text)
}

type textLine struct {
	size size
	nbrs []noBrRun
}

func (t *Text) AddRun(run *TextRun) *Text {
	t.Runs = append(t.Runs, run)
	return t
}
func (t *Text) SetAlign(align TextAlign) *Text {
	t.Align = align
	return t
}

func NewText(runs ...*TextRun) *Text {
	t := &Text{
		Runs: runs,
	}
	t.flexItemCommon.init(t)
	return t
}

func (t *Text) drawContent(pdf *gopdf.GoPdf, r rect) (err error) {
	defer wrap(&err, "text.drawContent")

	lines, err := t.splitLines(pdf, r.w)
	if err != nil {
		return err
	}

	pdf.SetXY(r.x, r.y)
	for _, line := range lines {
		pdf.SetX(r.x) // TODO align
		for _, nbr := range line.nbrs {
			if err := nbr.draw(pdf); err != nil {
				return err
			}
		}
		pdf.Br(line.size.h)
	}

	return nil
}

func (t *Text) getContentSize(pdf *gopdf.GoPdf, contentBoxMax size) (s size, err error) {
	defer wrap(&err, "text.getContentSize")

	lines, err := t.splitLines(pdf, contentBoxMax.w)
	if err != nil {
		return size{}, err
	}

	for _, line := range lines {
		s.w = math.Max(s.w, line.size.w)
		s.h += line.size.h
	}

	return s, nil
}

// splitLinesは Runs を行ごとに区切り、 [][]TextRunを返します。
// 返されるスライスは行を表しており、その要素は行に含まれるRunです。
// 下記のルールが考慮されます
// [ ] Textの幅
// [v] Runに含まれる改行コード
// [ ] 禁則処理
// [ ]  - 連続する欧文文字と空白
// [ ]  - 句読点や約物
func (t *Text) splitLines(pdf *gopdf.GoPdf, widthLimit float64) ([]textLine, error) {
	if false {
		pdf.SplitText("", 10)
	}

	lines := []textLine{}
	for _, r := range t.Runs {
		for i, nbr := range r.splitWithNewline() {
			if len(lines) == 0 || i != 0 {
				lines = append(lines, textLine{})
			}

			for {
				line := &lines[len(lines)-1]

				nbr1, nbr2, err := nbr.splitWithWidth(pdf, widthLimit-line.size.w)
				if err != nil {
					return nil, err
				}

				s, err := nbr1.size(pdf)
				if err != nil {
					return nil, err
				}

				line.nbrs = append(line.nbrs, *nbr1)
				line.size.w += s.w
				line.size.h = math.Max(line.size.h, s.h)

				if nbr2 != nil {
					lines = append(lines, textLine{})
					nbr = *nbr2
				} else {
					break
				}
			}
		}
	}

	return lines, nil
}
