package password

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordIncorrect = errors.New("incorrect password")

type Bcrypt struct {
	pepper string
}

func NewBcrypt(pepper string) *Bcrypt {
	return &Bcrypt{
		pepper: pepper,
	}
}

func (b *Bcrypt) GetHashedPassword(password string) ([]byte, error) {
	if password == "" {
		return nil, errors.New("empty password")
	}

	pwBytes := []byte(password + b.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedBytes, nil
}

func (b *Bcrypt) CompareHashedPassword(found, given string) error {
	pwBytes := []byte(given + b.pepper)
	err := bcrypt.CompareHashAndPassword([]byte(found), pwBytes)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return ErrPasswordIncorrect
		default:
			return err
		}
	}

	return nil
}
