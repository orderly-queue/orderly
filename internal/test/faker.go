package test

import "github.com/brianvoe/gofakeit/v7"

func init() {
	gofakeit.Seed(0)
}

func Word() string {
	return gofakeit.Word()
}

func Sentence(words int) string {
	return gofakeit.Sentence(words)
}

func Email() string {
	return gofakeit.Email()
}
