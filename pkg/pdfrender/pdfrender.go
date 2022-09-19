package pdfrender

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
)

type PDFRender struct {
	r *wkhtmltopdf.PDFGenerator
}

func NewPDFRenderer(r *wkhtmltopdf.PDFGenerator) *PDFRender {
	return &PDFRender{
		r: r,
	}
}

func (p PDFRender) RenderHTML(r io.Reader) error {
	return nil
}
