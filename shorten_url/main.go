package main

import (
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"

    "shorten_url/config"
    "shorten_url/helpers"
)

// Struct to represent a url
type Url struct {
    Url string
}

func returnOK(w http.ResponseWriter, value string){
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    m := map[string]string{
        "Status": "OK",
        "Message": value,
    }
    json.NewEncoder(w).Encode(m)
}

func returnERR(w http.ResponseWriter, message string, err error){
    log.Println(message, err)
    m := map[string]string{
        "Status": "Error",
        "Message": message,
    }
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}

func ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {

    var body Url
    var BaseUrlConf config.BaseUrlConf
    config.GetEnv(&BaseUrlConf)

    baseShorterUrl := BaseUrlConf.Protocol + "://" + BaseUrlConf.Url + "/"

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        returnERR(w, "Failed to decode json", err)
        return
    }
    if body.Url == "" {
        returnERR(w, "URL can't be empty", nil)
        return
    }

    // Try to get short URL value from Redis - return if it exists
    ok, val := helpers.GetFromRedis(body.Url)
    if !ok {
        returnERR(w, val, nil)
        return
    }

    // Value is not on redis - try to fetch from DB
    if val == "" {
        ok, val = helpers.GetFromPostgres("short", "url_to_short", "real_url", body.Url)
        if !ok {
            returnERR(w, val, nil)
            return
        }
    }
    if val != "" {
        log.Println("Url already exists, returning its value!")
        returnOK(w, baseShorterUrl + val)
    } else {
        // Generate short URL
        ok, shorterValue := helpers.ShortenUrl(body.Url)
        if !ok {
            returnERR(w, "Error while trying to hash", nil)
            return
        }

        // Insert result to DB
        ok = helpers.InsertToRedis(body.Url, shorterValue)
        ok2 := helpers.InsertToRedis(shorterValue, body.Url)
        if !ok || !ok2 {
            returnERR(w, "Error while trying to insert value to redis DB", nil)
            return
        }

        ok = helpers.InsertToPostgres("url_to_short", [2]string{"real_url", "short"}, [2]string{body.Url, shorterValue})
        ok2 = helpers.InsertToPostgres("short_to_url", [2]string{"short", "real_url"}, [2]string{shorterValue, body.Url})
        if !ok {
            returnERR(w, "Error while trying to insert value to postgres DB", nil)
        } else {
            returnOK(w, baseShorterUrl + shorterValue)
        }
    }

    
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/shorten_url", ShortenUrlHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":5000", router))
}