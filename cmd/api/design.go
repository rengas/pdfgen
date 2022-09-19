package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rengas/pdfgen/pkg/design"
	"github.com/rengas/pdfgen/pkg/httputils"
	"html/template"

	"net/http"
)

type DesignAPI struct {
	designRepo DesignRepository
}

func NewDesignAPI(designRepo DesignRepository) *DesignAPI {
	return &DesignAPI{
		designRepo: designRepo,
	}
}

type Field struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type CreateTemplateRequest struct {
	Name      string  `json:"name"`
	ProfileId string  `json:"profileId"`
	Design    string  `json:"design"`
	Fields    []Field `json:"fields"`
}

func (c CreateTemplateRequest) GetDesignModel() (design.Design, error) {

	if c.Name == "" {
		return design.Design{}, errors.New("name is empty")
	}

	if c.ProfileId == "" {
		return design.Design{}, errors.New("profile is empty")
	}

	if c.Fields == nil {
		return design.Design{}, errors.New("fields is empty")
	}

	for _, field := range c.Fields {
		f := field.Value
		switch f.(type) {
		case string:
		case float64:
		case int:
			continue
		default:
			return design.Design{}, errors.New(fmt.Sprintf("%s should either be float,string or int", field.Name))
		}
	}

	b, err := json.Marshal(c.Fields)
	if err != nil {
		return design.Design{}, errors.New("invalid fields structure")
	}

	if c.Design == "" {
		return design.Design{}, errors.New("design cannot be empty")
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(c.Design)
	if err != nil {
		return design.Design{}, errors.New("template must be base64 encoded")
	}

	//validate if valid design
	_, err = template.New("").Parse(string(rawDecodedText))
	if err != nil {
		return design.Design{}, fmt.Errorf("invalid html design %w", err)
	}

	return design.Design{
		Id:        uuid.NewString(),
		ProfileId: c.ProfileId,
		Name:      c.Name,
		Template:  rawDecodedText,
		Fields:    b,
	}, nil
}

func (d *DesignAPI) CreateDesign(w http.ResponseWriter, req *http.Request) {
	var t CreateTemplateRequest
	err := httputils.ReadJson(req, &t)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	dm, err := t.GetDesignModel()
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	err = d.designRepo.Save(context.Background(), dm)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to save profile"),
			http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, httputils.OKResponse{Id: dm.Id}, http.StatusOK)
}
