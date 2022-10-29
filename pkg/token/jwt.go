package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var ErrorUserIdEmpty = errors.New("userId is empty")

type JWT struct {
	accessSecretKey  string
	accessExpires    int
	refreshSecretKey string
	refreshExpires   int
	signingMethod    *jwt.SigningMethodHMAC
}

// TokenDetails - login response with access/refresh token pair.
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWT(accessSecretKey, refreshSecretKey string, accessExpires, refreshExpires int) *JWT {
	return &JWT{
		accessSecretKey:  accessSecretKey,
		accessExpires:    accessExpires,
		refreshSecretKey: refreshSecretKey,
		refreshExpires:   refreshExpires,
		signingMethod:    jwt.SigningMethodHS256,
	}
}

func (j JWT) TokePair(claims map[string]interface{}) (TokenDetails, error) {

	id, ok := claims["userId"]
	if !ok {
		return TokenDetails{}, ErrorUserIdEmpty
	}
	td := TokenDetails{}

	aClaims := jwt.MapClaims{}

	aClaims["exp"] = time.Now().Add(time.Minute * time.Duration(j.accessExpires)).Unix()
	aClaims["sub"] = id

	aToken := jwt.NewWithClaims(j.signingMethod, aClaims)

	accessToken, err := aToken.SignedString([]byte(j.accessSecretKey))
	if err != nil {
		return TokenDetails{}, err
	}
	td.AccessToken = accessToken

	rClaims := jwt.MapClaims{}

	rClaims["exp"] = time.Now().Add(time.Hour * time.Duration(j.refreshExpires)).Unix()
	rClaims["sub"] = id

	rToken := jwt.NewWithClaims(j.signingMethod, rClaims)
	refreshToken, err := rToken.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		return TokenDetails{}, err
	}

	td.RefreshToken = refreshToken

	return td, nil
}

func (j JWT) verifyToken(tkn string) (*jwt.Token, error) {
	token, err := jwt.Parse(tkn,
		func(token *jwt.Token) (interface{}, error) {
			//TODO what should you do here?
			return []byte(j.accessSecretKey), nil
		}, jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j JWT) verifyRefreshToken(tkn string) (*jwt.Token, error) {
	token, err := jwt.Parse(tkn,
		func(token *jwt.Token) (interface{}, error) {
			//TODO what should you do here?
			return []byte(j.refreshSecretKey), nil
		}, jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

type Claims map[string]interface{}

func (j JWT) ExtractTokenMetadata(tkn string) (Claims, error) {
	t, err := j.verifyToken(tkn)
	if err != nil {
		return Claims{}, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, nil
	}
	c := make(Claims, 0)
	c["userId"] = claims["sub"].(string)
	return c, nil
}
