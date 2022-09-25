package main

import (
	"context"
	"errors"
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
	Email      string `json:"email"`
	Provider   string `json:"provider"`
	FirebaseId string `json:"firebaseId"`
}

func (c CreateProfileRequest) Validate() error {

	if c.Email == "" {
		return errors.New("email is empty")
	}

	if c.Provider == "" {
		return errors.New("provider is empty")
	}

	if c.FirebaseId == "" {
		return errors.New("provider is empty")
	}

	return nil
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

	if err = pr.Validate(); err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest(err.Error()),
			http.StatusBadRequest)
		return
	}

	acc := account.Profile{
		Id:         uuid.NewString(),
		Email:      pr.Email,
		Provider:   pr.Provider,
		FirebaseId: pr.FirebaseId,
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

func (p *ProfileAPI) GetProfile(w http.ResponseWriter, req *http.Request) {
	profileId, ok := req.Context().Value("profileId").(string)
	if !ok {
		httputils.WriteJSON(w,
			httputils.BadRequest("profile not found"),
			http.StatusBadRequest)
	}

	prof, err := p.profileRepo.GetById(context.TODO(), profileId)
	if err != nil {
		httputils.WriteJSON(w,
			httputils.BadRequest("Unable get profile data "),
			http.StatusBadRequest)
		return
	}
	httputils.WriteJSON(w, prof, http.StatusOK)
}

func (p *ProfileAPI) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I'm ok")
}
