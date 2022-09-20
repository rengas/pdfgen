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
	m.Add(mediaTypeHTML, &html.Minifier{
		KeepConditionalComments: true,
		KeepDefaultAttrVals:     true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
		KeepWhitespace:          false,
	})
	return &Minifier{
		m: m,
	}
}

func (f Minifier) HTML(s string) (string, error) {
	bs, err := f.m.String(mediaTypeHTML, s)
	if err != nil {
		return "", ErrUnableToMinify
	}
	return bs, nil
}
