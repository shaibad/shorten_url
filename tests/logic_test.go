package tests


import (
	"testing"

	"github.com/stretchr/testify/assert"

	"url-shortener/helpers"
)

var (
	Url = "https://www.test.com"
	Shorter string
)

func TestShorten(t *testing.T) {
	var err error
	Shorter, err = helpers.ShortenUrl(Url)

	assert.NoError(t, err)
	assert.Equal(t, len([]rune(Shorter)), int(7))
}

func TestShortenUrlEmpty(t *testing.T) {
	var err error

	Shorter, err = helpers.ShortenUrl("")
	assert.Error(t, err)
}
