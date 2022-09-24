package middleware

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"net/http"
	"strings"
)

type FireBaseAuth interface {
	Verify(ctx context.Context, idToken string) (*auth.Token, error)
}

type Auth struct {
	f FireBaseAuth
}

func NewFirebaseAuth(f FireBaseAuth) *Auth {
	return &Auth{
		f: f,
	}
}

func (a Auth) FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		_, err := a.f.Verify(context.TODO(), idToken)
		if err != nil {
			//TODO What should be the header here?
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, ""))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
