package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"io"
	"log"
)

type AuthClient struct {
	a *auth.Client
}

func New(r io.Reader) *AuthClient {

	buf, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	opt := option.WithCredentialsJSON(buf)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return &AuthClient{
		a: a,
	}
}

func (f AuthClient) Verify(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := f.a.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
