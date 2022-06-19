package helpers

import (
	"log"
	"golang.org/x/crypto/bcrypt"
    "github.com/jcoene/go-base62"
	"math/big"
)

func ShortenUrl(url string) (bool, string) {
	valueToHash := []byte(url)
	// Hash original URL
	hash, err := bcrypt.GenerateFromPassword(valueToHash, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return false, "Error while trying to hash"
	}
	// Convert to Base 62 to allow correct url representation
	generatedNumber := new(big.Int).SetBytes(hash).Int64()
	shorterValue := base62.Encode(generatedNumber)[0:7]
	return true, shorterValue
}