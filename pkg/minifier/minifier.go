package minifier

import (
	"errors"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

const mediaTypeHTML = "text/html"

var ErrUnableToMinify = errors.New("unable to minifier")

type Minifier struct {
	m *minify.M
}

func NewMinifier() *Minifier {
	m := minify.New()
	m.AddFunc(mediaTypeHTML, html.Minify)
	return &Minifier{
		m: m,
	}
}

func (f Minifier) HTML(b []byte) ([]byte, error) {
	bs, err := f.m.Bytes(mediaTypeHTML, b)
	if err != nil {
		return nil, ErrUnableToMinify
	}
	return bs, nil
}
