package flexpdf

import (
	_ "embed"
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

	box := &Box{
		Items: []FlexItem{
			&Text{Text: "あいうえお"},
			&Text{Text: "かきくけこ"},
		},
	}

	if err := Draw(pdf, box, gopdf.PageSizeA4); err != nil {
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
