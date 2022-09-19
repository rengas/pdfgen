package main

import (
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
	"log"
)

type RenderRepository interface {
	RenderHTML(r io.Reader) error
}

func main() {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage("https://www.google.com.sg/?gfe_rd=cr&ei=c1kkWMnRC4uDvATQyaHgCw")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done

}
