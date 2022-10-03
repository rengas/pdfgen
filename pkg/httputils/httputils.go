package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rengas/pdfgen/pkg/design"
	"net/http"
)

type OkResponse struct {
	Id      string `json:"Id"`
	Message string `json:"Message"`
}
type ErrorResponse struct {
	Error error
}

func OK(ctx context.Context, w http.ResponseWriter, v interface{}) {
	WriteJSON(ctx, w, v, http.StatusOK)
}

func UnProcessableEntity(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusUnprocessableEntity)
}

func BadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusBadRequest)
}

func UnAuthorized(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusUnauthorized)
}

func Forbidden(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusForbidden)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusInternalServerError)
}

func NotFound(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusNotFound)
}

func Conflict(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusConflict)
}

func ReadJson(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, v interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("X-Request-ID", middleware.GetReqID(ctx))
	w.WriteHeader(statusCode)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func WriteFile(w http.ResponseWriter, b []byte, statusCode int) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(statusCode)
	w.Write(b)
	return nil
}

func WritePaginatedJSON(w http.ResponseWriter, pagination design.Pagination, v interface{}, statusCode int) error {
	w.Header().Add("count", fmt.Sprintf("%d", pagination.Page))
	w.Header().Add("total", fmt.Sprintf("%d", pagination.Total))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}
