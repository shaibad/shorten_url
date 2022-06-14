package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "github.com/jcoene/go-base62"
    "math/big"
)

type Url struct {
    Url string
}

var longToShortMap = make(map[string]string)
var shortToLongMap = make(map[string]string)

func ShortenUrl(w http.ResponseWriter, r *http.Request){

    var body Url
    err := json.NewDecoder(r.Body).Decode(&body)

    if err != nil {
        log.Println(err)
    }

    if _, ok := longToShortMap[body.Url]; ok {
        log.Println("Error, url already exists!")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Error, url already exists!"))
        return
    }

    valueToHash := []byte(body.Url)
    hash, err := bcrypt.GenerateFromPassword(valueToHash, bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error while trying to hash", err)
    }

    generatedNumber := new(big.Int).SetBytes(hash).Int64()
    shorterValue := base62.Encode(generatedNumber)

    longToShortMap[body.Url] = shorterValue
    shortToLongMap[shorterValue] = body.Url // To use in Get request

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Suceess"))
}

func GetUrl(w http.ResponseWriter, r *http.Request){
    shortUrlBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println("Error while trying to hash", err)
    }
    log.Println(string(shortUrlBytes))
    
    shortUrlStr := string(shortUrlBytes)
    if _, ok := longToShortMap[shortUrlStr]; ok {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(longToShortMap[shortUrlStr]))
    } else {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error, url not found!"))
    }
}

func handleRequests() {
	
    myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/shorten_url", ShortenUrl).Methods("POST")
    myRouter.HandleFunc("/get_url", GetUrl).Methods("GET")

    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    fmt.Println("URL shortner app")
    handleRequests()
}