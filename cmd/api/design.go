package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rengas/pdfgen/pkg/design"
	"github.com/rengas/pdfgen/pkg/httputils"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
)

type DesignAPI struct {
	designRepo DesignRepository
	minfier    Minifier
}

func NewDesignAPI(designRepo DesignRepository,
	minifier Minifier) *DesignAPI {
	return &DesignAPI{
		designRepo: designRepo,
		minfier:    minifier,
	}
}

type CreateTemplateRequest struct {
	Name      string       `json:"name"`
	ProfileId string       `json:"profileId"`
	Design    string       `json:"design"`
	Fields    design.Attrs `json:"fields"`
}

func (c CreateTemplateRequest) Validate() error {

	if c.Name == "" {
		return errors.New("name is empty")
	}

	if c.ProfileId == "" {
		return errors.New("profileId is empty")
	}

	if c.Design == "" {
		return errors.New("design is empty")
	}

	if c.Fields != nil {
		for k, v := range c.Fields {
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

func (d *DesignAPI) CreateDesign(w http.ResponseWriter, req *http.Request) {
	var t CreateTemplateRequest
	err := httputils.ReadJson(req, &t)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	err = t.Validate()
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	ds := design.Design{
		Id:        uuid.NewString(),
		ProfileId: t.ProfileId,
		Name:      t.Name,
		Fields:    nil,
	}

	dt, err := base64.StdEncoding.DecodeString(t.Design)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("design must be base64 encoded"),
			http.StatusBadRequest)
		return
	}

	ws := string(dt)
	//validate if valid design
	_, err = template.New(t.Name).Parse(ws)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("invalid html design "),
			http.StatusBadRequest)
		return
	}

	mt, err := d.minfier.HTML(ws)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("unable to minifier template"),
			http.StatusInternalServerError)
		return
	}

	ds.Template = mt

	if t.Fields != nil {
		ds.Fields = &t.Fields
	}

	err = d.designRepo.Save(context.Background(), ds)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to update design"),
			http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, httputils.OKResponse{Id: ds.Id}, http.StatusOK)
}

type UpdateTemplateRequest struct {
	Name   string       `json:"name"`
	Design string       `json:"design"`
	Fields design.Attrs `json:"fields"`
}

func (c UpdateTemplateRequest) Validate() error {

	if c.Name == "" {
		return errors.New("name is empty")
	}

	if c.Design == "" {
		return errors.New("design is empty")
	}

	if c.Fields != nil {
		for k, v := range c.Fields {
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

func (d *DesignAPI) UpdateDesign(w http.ResponseWriter, req *http.Request) {
	designId := chi.URLParam(req, "designId")
	if designId == "" {
		httputils.WriteJSON(w,
			httputils.BadRequest("designId is empty"),
			http.StatusBadRequest)
	}

	var t UpdateTemplateRequest
	err := httputils.ReadJson(req, &t)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	err = t.Validate()
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	ds := design.Design{
		Id:     designId,
		Name:   t.Name,
		Fields: nil,
	}

	dt, err := base64.StdEncoding.DecodeString(t.Design)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("design must be base64 encoded"),
			http.StatusBadRequest)
		return
	}

	ws := string(dt)
	//validate if valid design
	_, err = template.New(t.Name).Parse(ws)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("invalid html design"),
			http.StatusBadRequest)
		return
	}

	mt, err := d.minfier.HTML(ws)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("unable to minifier template"),
			http.StatusInternalServerError)
		return
	}

	ds.Template = mt

	if t.Fields != nil {
		ds.Fields = &t.Fields
	}

	err = d.designRepo.Update(context.Background(), ds)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to update design"),
			http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, httputils.OKResponse{Id: ds.Id}, http.StatusOK)
}

func (d *DesignAPI) GetDesign(w http.ResponseWriter, req *http.Request) {
	designId := chi.URLParam(req, "designId")
	if designId == "" {
		httputils.WriteJSON(w,
			httputils.BadRequest("designId is empty"),
			http.StatusBadRequest)
	}

	ds, err := d.designRepo.GetById(context.TODO(), designId)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to get design by id"),
			http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, ds, http.StatusOK)
}

func (d *DesignAPI) ListDesign(w http.ResponseWriter, req *http.Request) {

	profileId, ok := req.Context().Value("profileId").(string)
	if !ok {
		httputils.WriteJSON(w,
			httputils.BadRequest("profileId is empty"),
			http.StatusBadRequest)
	}

	count := req.URL.Query().Get("count")
	if count == "" {
		httputils.WriteJSON(w,
			httputils.BadRequest("count is empty"),
			http.StatusBadRequest)
	}

	c, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("invalid count"),
			http.StatusBadRequest)
	}

	page := req.URL.Query().Get("page")
	if page == "" {
		httputils.WriteJSON(w,
			httputils.BadRequest("page is empty"),
			http.StatusBadRequest)
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("invalid page"),
			http.StatusBadRequest)
	}

	q := req.URL.Query().Get("search")
	lq := design.ListQuery{
		ProfileId: profileId,
		Limit:     c,
		Page:      p,
	}

	var ds []design.Design
	var pagi design.Pagination

	if q != "" {
		ds, pagi, err = d.designRepo.Search(context.TODO(), lq)

	} else {
		ds, pagi, err = d.designRepo.ListByProfileId(context.TODO(), lq)
	}

	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("unable to get List of designs"),
			http.StatusBadRequest)
	}

	httputils.WritePaginatedJSON(w,
		pagi,
		ds,
		http.StatusOK)
}

func (d *DesignAPI) DeleteDesign(w http.ResponseWriter, req *http.Request) {
	designId := chi.URLParam(req, "designId")
	if designId == "" {
		httputils.WriteJSON(w,
			httputils.BadRequest("designId is empty"),
			http.StatusBadRequest)
	}

	err := d.designRepo.Delete(context.TODO(), designId)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to get design by id"),
			http.StatusInternalServerError)
		return
	}
	httputils.WriteJSON(w, httputils.OKResponse{Msg: "design was deleted"}, http.StatusOK)
}
