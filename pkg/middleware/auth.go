package middleware

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/rengas/pdfgen/pkg/account"
	"log"
	"net/http"
	"strings"
)

type FireBaseAuth interface {
	Verify(ctx context.Context, idToken string) (*auth.Token, error)
}
type ProfileRepository interface {
	GetByFirebaseId(ctx context.Context, id string) (account.Profile, error)
}

type Auth struct {
	f FireBaseAuth
	p ProfileRepository
}

func NewFirebaseAuth(f FireBaseAuth, p ProfileRepository) *Auth {
	return &Auth{
		f: f,
		p: p,
	}
}

func (a Auth) FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		token, err := a.f.Verify(context.TODO(), idToken)
		if err != nil {
			//TODO What should be the header here?
			log.Printf("%s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		prof, err := a.p.GetByFirebaseId(context.TODO(), token.UID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "profileId", prof.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a Auth) FirebaseAuthWithout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		_, err := a.f.Verify(context.TODO(), idToken)
		if err != nil {
			//TODO What should be the header here?
			log.Printf("%s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
