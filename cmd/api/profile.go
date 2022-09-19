package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rengas/pdfgen/pkg/account"
	"github.com/rengas/pdfgen/pkg/httputils"
	"net/http"
)

type ProfileAPI struct {
	profileRepo ProfileRepository
}

func NewProfileAPI(profileRepo ProfileRepository) *ProfileAPI {
	return &ProfileAPI{
		profileRepo: profileRepo,
	}
}

type CreateProfileRequest struct {
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
}

func (p *ProfileAPI) CreateProfile(w http.ResponseWriter, req *http.Request) {
	var pr CreateProfileRequest
	err := httputils.ReadJson(req, &pr)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("Unable to read request"),
			http.StatusBadRequest)
		return
	}

	acc := account.Profile{
		Id:    uuid.NewString(),
		Email: pr.Email,
	}
	err = p.profileRepo.Save(context.Background(), acc)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.InternalError("Unable to save profile"),
			http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, httputils.OKResponse{Id: acc.Id}, http.StatusOK)
}

func (p *ProfileAPI) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I'm ok")
}
