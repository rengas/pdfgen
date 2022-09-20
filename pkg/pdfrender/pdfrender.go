package pdfrender

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
	"log"
)

type PDFRender struct {
	r *wkhtmltopdf.PDFGenerator
}

func NewPDFRenderer() *PDFRender {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	return &PDFRender{
		r: pdfg,
	}
}

func (p PDFRender) HTML(r io.Reader) ([]byte, error) {

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPageReader(r)

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	p.r.AddPage(page)

	// Create PDF document in internal buffer
	err := p.r.Create()
	if err != nil {
		return nil, err
	}

	// Output: Done
	return p.r.Bytes(), nil
}
