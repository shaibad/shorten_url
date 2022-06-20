package helpers

import (
	"log"
	"golang.org/x/crypto/bcrypt"
    "github.com/jcoene/go-base62"
	"math/big"
	"net/http"
	"encoding/json"
)

func ReturnOK(w http.ResponseWriter, value string){
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    m := map[string]string{
        "Status": "OK",
        "Message": value,
    }
    json.NewEncoder(w).Encode(m)
}

func ReturnERR(w http.ResponseWriter, message string, err error){
    log.Println(message, err)
    m := map[string]string{
        "Status": "Error",
        "Message": message,
    }
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}

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