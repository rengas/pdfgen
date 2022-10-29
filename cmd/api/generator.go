package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"html/template"
	"net/http"
	"strings"
)

type GeneratorAPI struct {
	designRepo DesignRepository
	renderer   Renderer
}

func NewGeneratorAPI(designRepo DesignRepository, renderer Renderer) *GeneratorAPI {
	return &GeneratorAPI{
		designRepo: designRepo,
		renderer:   renderer,
	}
}

// GeneratePDF func for updating new design.
// @Description  Generate a pdf
// @Summary      GeneratePDF
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param        GeneratePDFRequest body  GeneratePDFRequest  true  "register details"
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design/generate [post]
func (d *GeneratorAPI) GeneratePDF(w http.ResponseWriter, req *http.Request) {
	var t GeneratePDFRequest
	err := httputils.ReadJson(req, &t)
	if err != nil {
		httputils.BadRequest(context.TODO(), w, errors.New("unable to read request"))
		return
	}

	err = t.Validate()
	if err != nil {
		httputils.BadRequest(context.TODO(), w, err)
		return
	}

	design, err := d.designRepo.GetById(context.TODO(), t.DesignId)
	if err != nil {
		httputils.BadRequest(context.TODO(), w, errors.New("unable to get design"))
		return
	}

	if design.Fields != nil {
		tl, err := template.New(design.Name).Parse(design.Template)
		if err != nil {
			httputils.BadRequest(context.TODO(), w, errors.New("unable to parse template"))
			return
		}

		var buf bytes.Buffer

		err = tl.Execute(&buf, design.Fields)
		if err != nil {
			httputils.BadRequest(context.TODO(), w, errors.New("unable to match fields to design"))
			return
		}

		pb, err := d.renderer.HTML(&buf)
		if err != nil {
			httputils.BadRequest(context.TODO(), w, errors.New("unable to renderer pdf"))
			return
		}

		httputils.WriteFile(w,
			pb,
			http.StatusOK)
		return
	}

	pb, err := d.renderer.HTML(strings.NewReader(design.Template))
	if err != nil {
		httputils.InternalServerError(context.TODO(), w, errors.New("unable to renderer pdf"))
		return
	}

	httputils.WriteFile(w,
		pb,
		http.StatusOK)
	return

}
