package flexpdf

import (
	_ "embed"
	"image/color"
	"os"
	"testing"

	"github.com/signintech/gopdf"
)

var (
	//go:embed "testdata/fonts/ipaexg.ttf"
	ipaexgBytes []byte
	//go:embed "testdata/fonts/ipaexm.ttf"
	ipaexmBytes []byte
)

func TestXxx(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})

	if err := pdf.AddTTFFontData("ipaexg", ipaexgBytes); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFontData("ipaexm", ipaexmBytes); err != nil {
		t.Fatal(err)
	}

	root := &Box{
		Direction: DirectionColumn,
		Items: []FlexItem{
			&Box{
				BackgroundColor: color.RGBA{G: 0xFF, A: 0x80},
				Items: []FlexItem{
					&Box{},
					&Text{Text: "あいうえお\nかきくけこさしすせそ\nたちつ", FontFamily: "ipaexg", FontSize: 20, LineHeight: 1.5},
					&Text{Text: "abc", FontFamily: "ipaexm", FontSize: 30},
				},
			},
			&Text{Text: "あいうえお\nかきくけこさしすせそ\nたちつ", FontFamily: "ipaexg", FontSize: 20, LineHeight: 1.5},
			&Text{Text: "abc", FontFamily: "ipaexm", FontSize: 30},
		},
	}

	if err := Draw(pdf, root, gopdf.PageSizeA4); err != nil {
		t.Fatal(err)
	}

	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("out.pdf", data, 0666); err != nil {
		t.Fatal(err)
	}
}
