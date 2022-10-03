package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rengas/pdfgen/pkg/design"
	"github.com/rengas/pdfgen/pkg/httputils"
	"html/template"
	"net/http"
	"reflect"
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

type GeneratePDFRequest struct {
	DesignId string       `json:"DesignId"`
	Fields   design.Attrs `json:"fields"`
}

func (g GeneratePDFRequest) Validate() error {
	if g.DesignId == "" {
		return errors.New("designId is empty")
	}

	if g.Fields != nil {
		for k, v := range g.Fields {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Slice,
				reflect.Array, reflect.Map:
				continue
			default:
				return errors.New(fmt.Sprintf("%s has unsupported type for value", k))
			}
		}
	}

	return nil
}

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
