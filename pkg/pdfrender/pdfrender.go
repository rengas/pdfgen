package pdfrender

import (
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
	"log"
)

type PDFRender struct {
}

func NewPDFRenderer() *PDFRender {
	return &PDFRender{}
}

func (p PDFRender) HTML(r io.Reader) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPageReader(r)
	page.EnableLocalFileAccess.Set(true)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Output: Done
	return pdfg.Bytes(), nil
}
