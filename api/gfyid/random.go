package gfyid

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
	"unicode"
)

type jsonContent struct {
	Adjectives []string `json:"adjectives"`
	Nouns      []string `json:"nouns"`
}

var jsonData = jsonContent{}
var privateRand *rand.Rand

func init() {
	privateRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	jsonData = jsonContent{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		log.Fatal(err)
	}
}

// Returns a random part of a slice
func randomFrom(source []string) string {
	return source[privateRand.Intn(len(source))]
}

func upper(word string) string {
	a := []rune(word)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

// Noun returns a random noun
func Noun() string {
	return randomFrom(jsonData.Nouns)
}

// Adjective returns a random adjective
func Adjective() string {
	return randomFrom(jsonData.Adjectives)
}

// RandomID returns a random ID in the form of 'AdjectiveAdjectiveNoun`
func RandomID() string {
	return upper(Adjective()) + upper(Adjective()) + upper(Noun())
}
