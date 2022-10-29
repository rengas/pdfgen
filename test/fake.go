package test

import "github.com/brianvoe/gofakeit/v6"

func int64Ptr(n int64) *int64 {
	return &n
}

func float64Ptr(n float64) *float64 {
	return &n
}

func fromStringPtr(n *string) string {
	return *n
}
func StringPtr(n string) *string {
	return &n
}

func boolPtr(b bool) *bool {
	return &b
}

func Name() *string {
	name := gofakeit.Name()
	return &name
}

func LoremIpsumParagraph(paragraphCount int, sentenceCount int, wordCount int, separator string) *string {
	para := gofakeit.LoremIpsumParagraph(paragraphCount, sentenceCount, wordCount, separator)
	return &para
}

func StreetName() *string {
	para := gofakeit.StreetName()
	return &para
}

func Latitude() *float64 {
	para := gofakeit.Latitude()
	return &para
}

func Longitude() *float64 {
	para := gofakeit.Longitude()
	return &para
}

func PostalCode() *string {
	para := gofakeit.Zip()
	return &para
}

func Int64Range(min, max int64) int64 {
	return int64(gofakeit.IntRange(int(min), int(max)))
}

func IntRange(min, max int) int {
	return gofakeit.IntRange(min, max)
}

func Email() string {
	return gofakeit.Email()
}

func Password(lower bool, upper bool, numeric bool, special bool, space bool, num int) string {
	return gofakeit.Password(lower, upper, numeric, special, space, num)
}
func Phone() string {
	return gofakeit.Phone()
}

func ImageJpeg() []byte {
	return gofakeit.ImageJpeg(200, 200)
}

func ImagePng() []byte {
	return gofakeit.ImagePng(IntRange(200, 400), IntRange(200, 400))
}
