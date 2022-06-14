package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    // "io/ioutil"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "github.com/jcoene/go-base62"
    "math/big"
)

// Struct to represent a url
type Url struct {
    Url string
}

// Store maps between short and long urls
var longToShortMap = make(map[string]string)
var shortToLongMap = make(map[string]string)

func ShortenUrl(w http.ResponseWriter, r *http.Request){

    var body Url
    err := json.NewDecoder(r.Body).Decode(&body)

    if err != nil {
        log.Println(err)
    }

    if _, ok := longToShortMap[body.Url]; ok {
        log.Println("Url already exists, returning its value!")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(longToShortMap[body.Url]))
        return
    }

    valueToHash := []byte(body.Url)

    // Hash original URL
    hash, err := bcrypt.GenerateFromPassword(valueToHash, bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error while trying to hash", err)
    }

    // Convert to Base 62 to allow correct url representation
    generatedNumber := new(big.Int).SetBytes(hash).Int64()
    shorterValue := base62.Encode(generatedNumber)

    longToShortMap[body.Url] = shorterValue
    shortToLongMap[shorterValue] = body.Url // To use in Get request

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(shorterValue))
}

func GetUrl(w http.ResponseWriter, r *http.Request){
    // shortUrlBytes, _ := ioutil.ReadAll(r.Body)
    // shortUrl, err := ioutil.ReadAll(r.URL)

    url := r.URL.Query()["url"][0]
    log.Println(url)
    // log.Println(string(shortUrlBytes))
    
    shortUrlStr := string(url)
    // Check if mapping exists
    if _, ok := shortToLongMap[shortUrlStr]; ok {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(shortToLongMap[shortUrlStr]))
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