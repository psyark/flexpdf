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
			testJustifyContent(JustifyContentFlexStart),
			testJustifyContent(JustifyContentFlexEnd),
			testJustifyContent(JustifyContentCenter),
			testJustifyContent(JustifyContentSpaceBetween),
			testJustifyContent(JustifyContentSpaceAround),
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

func testJustifyContent(jc JustifyContent) *Box {
	return &Box{
		Direction: DirectionColumn,
		Items: []FlexItem{
			&Text{Text: string(jc) + ":", FontFamily: "ipaexg", FontSize: 20},
			&Box{
				BackgroundColor: color.RGBA{0x88, 0x88, 0x88, 0xFF},
				JustifyContent:  jc,
				Items: []FlexItem{
					&Box{
						BackgroundColor: color.RGBA{0xFF, 0xCC, 0xCC, 0xFF},
						Items:           []FlexItem{&Text{Text: "あいうえお", FontFamily: "ipaexg", FontSize: 24}},
					},
					&Box{
						BackgroundColor: color.RGBA{0xCC, 0xFF, 0xCC, 0xFF},
						Items:           []FlexItem{&Text{Text: "あいうえお", FontFamily: "ipaexg", FontSize: 24}},
					},
					&Box{
						BackgroundColor: color.RGBA{0xCC, 0xCC, 0xFF, 0xFF},
						Items:           []FlexItem{&Text{Text: "あいうえお かきくけこ", FontFamily: "ipaexg", FontSize: 24}},
					},
				},
			},
		},
	}
}
