package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/rengas/pdfgen/pkg/contexts"
	pgerror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"github.com/rengas/pdfgen/pkg/logging"
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
// @Router       /generate [post]
func (d *GeneratorAPI) GeneratePDF(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
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

	userId, err := contexts.UserIdFromContext(ctx)
	if err != nil {
		logging.Debug("unable to get userId from context")
		httputils.BadRequest(context.TODO(), w, pgerror.ErrUnableToGetUserIdFromContext)
		return
	}

	design, err := d.designRepo.GetById(ctx, userId, t.DesignId)
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
