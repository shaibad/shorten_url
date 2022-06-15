package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
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

func returnOK(w http.ResponseWriter, value string){
    w.WriteHeader(http.StatusOK)
    m := map[string]string{
        "Status": "OK",
        "Message": value,
    }
    json.NewEncoder(w).Encode(m)
}

func returnERR(w http.ResponseWriter, message string){
    m := map[string]string{
        "Status": "Error",
        "Message": message,
    }
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(m)
}

func ShortenUrl(w http.ResponseWriter, r *http.Request) {

    var body Url
    err := json.NewDecoder(r.Body).Decode(&body)

    if err != nil {
        log.Println(err)
    }

    if body.Url == "" {
        log.Println("Error, URL can't be empty")
        returnERR(w, "URL can't be empty")
        return
    }

    if _, ok := longToShortMap[body.Url]; ok {
        log.Println("Url already exists, returning its value!")
        returnOK(w, longToShortMap[body.Url])
        return
    }

    valueToHash := []byte(body.Url)

    // Hash original URL
    hash, err := bcrypt.GenerateFromPassword(valueToHash, bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error while trying to hash", err)
        returnERR(w, "Error while trying to hash")
        return
    }

    // Convert to Base 62 to allow correct url representation
    generatedNumber := new(big.Int).SetBytes(hash).Int64()
    shorterValue := base62.Encode(generatedNumber)

    longToShortMap[body.Url] = shorterValue
    shortToLongMap[shorterValue] = body.Url // To use in Get request

    returnOK(w, shorterValue)
}

func GetUrl(w http.ResponseWriter, r *http.Request){
    urlParam := r.URL.Query()["url"]
    if len(urlParam) == 0 {
        log.Println("Error, URL param can't be empty")
        returnERR(w, "URL param can't be empty")
        return
    }

    url := r.URL.Query()["url"][0]
    shortUrlStr := string(url)
    w.Header().Set("Content-Type", "application/json")
    // Check if mapping exists
    if _, ok := shortToLongMap[shortUrlStr]; ok {
        w.WriteHeader(http.StatusOK)
        var u Url
        u.Url = shortToLongMap[shortUrlStr]
        json.NewEncoder(w).Encode(u)
    } else {
        returnERR(w, "Url not found")
    }
}

func handleRequests() {
	
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/shorten_url", ShortenUrl).Methods("POST")
    router.HandleFunc("/get_url", GetUrl).Methods("GET")

    log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
    fmt.Println("URL shortner app")
    handleRequests()
}