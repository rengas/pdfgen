package httputils

import (
	"encoding/json"
	"fmt"
	"github.com/rengas/pdfgen/pkg/design"
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type OKResponse struct {
	Id  string `json:"id"`
	Msg string `json:"msg"`
}

func BadRequest(msg string) ErrorResponse {
	return ErrorResponse{
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

func InternalError(msg string) ErrorResponse {
	return ErrorResponse{
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
	}
}

func ReadJson(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, v interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
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
