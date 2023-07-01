package flexpdf

import "github.com/signintech/gopdf"

func Draw(pdf *gopdf.GoPdf, box *Box, pageSize *gopdf.Rect) error {
	pdf.AddPageWithOption(gopdf.PageOption{PageSize: pageSize})
	if err := pdf.SetFont("ipaexg", "", 20); err != nil {
		return err
	}
	return pdf.CellWithOption(&gopdf.Rect{W: 100, H: 100}, "HOGE", gopdf.CellOption{Border: gopdf.AllBorders})
}
