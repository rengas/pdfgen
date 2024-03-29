package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rengas/pdfgen/pkg/contexts"
	"github.com/rengas/pdfgen/pkg/design"
	pgerror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"github.com/rengas/pdfgen/pkg/logging"
	"github.com/rengas/pdfgen/pkg/pagination"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type DesignAPI struct {
	designRepo DesignRepository
	minifier   Minifier
}

func NewDesignAPI(designRepo DesignRepository,
	minifier Minifier) *DesignAPI {
	return &DesignAPI{
		designRepo: designRepo,
		minifier:   minifier,
	}
}

// CreateDesign func for register.
// @Description  Create a new Design.
// @Summary      Create Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param        CreateDesignRequest body  CreateDesignRequest  true  "register details"
// @Success      200           {object}  CreateDesignResponse 		"Created"
// @Failure      400           {object}  httputils.ErrorResponse    "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse    "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse    "Internal Server Error"
// @Router       /design [post]
func (d *DesignAPI) CreateDesign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var t CreateDesignRequest
	err := httputils.ReadJson(r, &t)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, errors.New("unable to read request"))
		return
	}

	t.UserId = d.getUserId(w, r)
	if t.UserId == "" {
		logging.Debug("unable to get userId from context")
		httputils.NotFound(ctx, w, pgerror.ErrUnableToGetUserIdFromContext)
		return
	}

	err = t.Validate()
	if err != nil {
		logging.WithField(logging.Field{Label: "validation", Value: err.Error()}).Debug("validation failed")
		httputils.BadRequest(ctx, w, err)
		return
	}

	ds := design.Design{
		Id:     uuid.NewString(),
		UserId: t.UserId,
		Name:   t.Name,
		Fields: nil,
	}

	dt, err := base64.StdEncoding.DecodeString(t.Design)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignMustBeBase64Encoded)
		return
	}

	ws := string(dt)
	//validate if valid design
	_, err = template.New(t.Name).Parse(ws)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignInvalidHTML)
		return
	}

	mt, err := d.minifier.HTML(ws)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignUnableToMinify)
		return
	}

	ds.Template = mt

	if t.Fields != nil {
		ds.Fields = &t.Fields
	}

	err = d.designRepo.Save(context.Background(), ds)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignUnableToUpdate)
		return
	}

	httputils.OK(context.TODO(), w, CreateDesignResponse{Id: ds.Id})
}

// UpdateDesign func for updating new design.
// @Description  Update a Design.
// @Summary      Update Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param   	 designId     path    string     true        "design id"
// @Param        UpdateDesignRequest body  UpdateDesignRequest  true  "register details"
// @Success      200           {object}  UpdateDesignResponse "Created"
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design/{designId} [put]
func (d *DesignAPI) UpdateDesign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, designId := d.getUserIdAndDesignId(w, r)
	if userId == "" || designId == "" {
		return
	}

	ds, err := d.designRepo.GetById(context.TODO(), userId, designId)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(context.TODO(), w, ErrDesignUnableToGetDesign)
		return
	}

	var t UpdateDesignRequest
	err = httputils.ReadJson(r, &t)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignUnableToReadRequest)
		return
	}

	err = t.Validate()
	if err != nil {
		logging.WithField(logging.Field{Label: "validation", Value: err.Error()}).Debug("validation failed")
		httputils.BadRequest(context.TODO(), w, err)
		return
	}

	ds = design.Design{
		Id:     designId,
		Name:   t.Name,
		Fields: nil,
	}

	dt, err := base64.StdEncoding.DecodeString(t.Design)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignMustBeBase64Encoded)
		return
	}

	//validate if valid design
	ws := string(dt)
	_, err = template.New(t.Name).Parse(ws)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(context.TODO(), w, ErrDesignInvalidHTML)
		return
	}

	mt, err := d.minifier.HTML(ws)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(context.TODO(), w, ErrDesignUnableToMinify)
		return
	}

	ds.Template = mt

	if t.Fields != nil {
		ds.Fields = &t.Fields
	}
	ds.UpdatedAt = time.Now().UTC()

	err = d.designRepo.Update(context.Background(), ds)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(context.TODO(), w, ErrDesignUnableToUpdate)
		return
	}

	httputils.OK(context.TODO(), w, UpdateDesignResponse{Id: ds.Id})
}

// GetDesign func for updating new design.
// @Description  Get a Design.
// @Summary      Get Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param   	 designId     path    string     true   "design id"
// @Success      200           {object}  GetDesignResponse
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design/{designId} [get]
func (d *DesignAPI) GetDesign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, designId := d.getUserIdAndDesignId(w, r)
	if userId == "" || designId == "" {
		return
	}

	ds, err := d.designRepo.GetById(context.TODO(), userId, designId)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(context.TODO(), w, ErrDesignUnableToGetDesign)
		return
	}

	httputils.OK(context.TODO(), w, GetDesignResponse(ds))
}

// ListDesign func for updating new design.
// @Description  List Designs.
// @Summary      List Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Success      200           {object}  ListDesignResponse
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design [get]
func (d *DesignAPI) ListDesign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count := r.URL.Query().Get("count")
	if count == "" {
		logging.Debug("count is empty")
		httputils.BadRequest(context.TODO(), w, ErrDesignCountIsEmpty)
		return
	}

	c, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		logging.Debug("unable to parse count")
		httputils.BadRequest(context.TODO(), w, ErrDesignCountInvalid)
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		logging.Debug("page is empty ")
		httputils.BadRequest(context.TODO(), w, ErrDesignPageIsEmpty)
		return
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		logging.WithContext(ctx).WithError(err).Debug("unable to parse page")
		httputils.BadRequest(context.TODO(), w, ErrDesignPageInvalid)
		return
	}

	userId := d.getUserId(w, r)
	if userId == "" {
		logging.Debug("unable to get userId from context")
		httputils.NotFound(ctx, w, pgerror.ErrUnableToGetUserIdFromContext)
		return
	}

	lq := design.ListQuery{
		UserId: userId,
		Limit:  c,
		Page:   p,
	}

	var ds []design.Design
	var pagi pagination.Pagination

	q := r.URL.Query().Get("search")
	if q != "" {
		lq.Query = q
		ds, pagi, err = d.designRepo.Search(context.TODO(), lq)
		if err != nil {
			logging.WithContext(ctx).Error(err.Error())
			httputils.BadRequest(context.TODO(), w, ErrDesignUnableToGetDesigns)
			return
		}

	} else {
		ds, pagi, err = d.designRepo.ListByUserId(context.TODO(), lq)
		if err != nil {
			logging.WithContext(ctx).Error(err.Error())
			httputils.BadRequest(context.TODO(), w, ErrDesignUnableToGetDesigns)
			return
		}
	}

	httputils.WriteJSON(ctx, w,
		ListDesignResponse{
			Designs:    ds,
			Pagination: pagi,
		},
		http.StatusOK)
}

// DeleteDesign func for updating new design.
// @Description  Delete a Design.
// @Summary      Delete Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param   	 designId     path    string     true   "design id"
// @Success      200           {object}  DeleteDesignResponse
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design/{designId} [delete]
func (d *DesignAPI) DeleteDesign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, designId := d.getUserIdAndDesignId(w, r)
	if userId == "" || designId == "" {
		return
	}

	_, err := d.designRepo.GetById(context.TODO(), userId, designId)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(context.TODO(), w, ErrDesignUnableToGetDesign)
		return
	}

	err = d.designRepo.Delete(context.TODO(), userId, designId)
	if err != nil {
		logging.WithContext(ctx).Error(err.Error())
		httputils.InternalServerError(context.TODO(), w, errors.New("unable to get design by id"))
		return
	}
	httputils.OK(context.TODO(), w, DeleteDesignResponse{Id: designId})

}

// ValidateDesign func for updating new design.
// @Description  Validate a Design.
// @Summary      Validate Design
// @Tags         Design
// @Accept       json
// @Produce      json
// @Param        ValidateDesignRequest body  ValidateDesignRequest  true  "register details"
// @Success      200           {object}  ValidateDesignResponse
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /design/validate [post]
func (d *DesignAPI) ValidateDesign(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var t ValidateDesignRequest
	err := httputils.ReadJson(req, &t)
	if err != nil {
		logging.WithContext(ctx).Error(err.Error())
		httputils.BadRequest(context.TODO(), w, ErrDesignUnableToReadRequest)
		return
	}

	err = t.Validate()
	if err != nil {
		logging.WithField(logging.Field{Label: "validation", Value: err.Error()}).Debug("validation failed")
		httputils.BadRequest(context.TODO(), w, err)
		return
	}

	dt, err := base64.StdEncoding.DecodeString(t.Design)
	if err != nil {
		logging.WithContext(ctx).Error(err.Error())
		httputils.BadRequest(context.TODO(), w, ErrDesignMustBeBase64Encoded)
		return
	}

	ws := string(dt)
	//validate if valid design
	_, err = template.New(t.Name).Parse(ws)
	if err != nil {
		logging.WithContext(ctx).Error(err.Error())
		httputils.BadRequest(context.TODO(), w, ErrDesignInvalidHTML)
		return
	}

	if t.Fields != nil {
		tl, err := template.New(t.Name).Parse(t.Design)
		if err != nil {
			logging.WithContext(ctx).Error(err.Error())
			httputils.BadRequest(context.TODO(), w, ErrDesignUnableToParseDesign)
			return
		}

		var buf bytes.Buffer

		err = tl.Execute(&buf, t.Fields)
		if err != nil {
			logging.WithContext(ctx).Error(err.Error())
			httputils.BadRequest(context.TODO(), w, ErrDesignUnableToMatchFieldsToDesign)
			return
		}

	}
	httputils.OK(context.TODO(), w, ValidateDesignResponse{Message: "design is good to go"})
}

// getUserIdAndDesignId get userId from context
func (d *DesignAPI) getUserId(w http.ResponseWriter, req *http.Request) string {
	ctx := req.Context()
	userId, err := contexts.UserIdFromContext(ctx)
	if err != nil {
		logging.Debug("unable to get userId from context")
		httputils.BadRequest(context.TODO(), w, pgerror.ErrUnableToGetUserIdFromContext)
		return ""
	}
	return userId
}

//getUserIdAndDesignId get both userId and designId from context
func (d *DesignAPI) getUserIdAndDesignId(w http.ResponseWriter, req *http.Request) (string, string) {
	ctx := req.Context()
	designId := chi.URLParam(req, "designId")
	if designId == "" {
		logging.WithContext(ctx).Debug("unable to get designId from context")
		httputils.BadRequest(context.TODO(), w, errors.New("designId is empty"))
		return "", ""
	}

	userId, err := contexts.UserIdFromContext(ctx)
	if err != nil {
		logging.Debug("unable to get userId from context")
		httputils.BadRequest(context.TODO(), w, pgerror.ErrUnableToGetUserIdFromContext)
		return "", ""
	}
	return designId, userId
}
