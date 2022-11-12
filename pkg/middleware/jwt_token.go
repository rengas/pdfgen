package middleware

import (
	"github.com/rengas/pdfgen/pkg/contexts"
	mkerror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"github.com/rengas/pdfgen/pkg/token"
	"net/http"
	"strings"
)

type JWTToken interface {
	ExtractTokenMetadata(tkn string) (token.Claims, error)
}

type JWTMiddleware struct {
	jwt JWTToken
}

func NewJWTToken(j JWTToken) *JWTMiddleware {
	return &JWTMiddleware{
		jwt: j,
	}
}

func (j JWTMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			httputils.UnAuthorized(ctx, w, mkerror.ErrUnAuthorizedError)
			return
		}
		bToken := strings.Split(bearer, " ")
		if len(bToken) != 2 {
			httputils.UnAuthorized(ctx, w, mkerror.ErrUnAuthorizedError)
			return
		}

		c, err := j.jwt.ExtractTokenMetadata(bToken[1])
		if err != nil {
			httputils.UnAuthorized(ctx, w, mkerror.ErrUnAuthorizedError)
			return
		}
		userId, ok := c["userId"].(string)
		if !ok {
			httputils.UnAuthorized(ctx, w, mkerror.ErrUnAuthorizedError)
			return
		}

		ctx = contexts.WithUserId(ctx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
