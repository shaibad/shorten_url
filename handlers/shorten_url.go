package handlers

import (
    "log"
    "encoding/json"
    "net/http"

    "url-shortener/config"
    "url-shortener/helpers"
    "url-shortener/db"
)

// Struct to represent a url
type Url struct {
    Url string
}

func ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {

    var body Url
    var BaseUrlConf config.BaseUrlConf
    config.GetEnv(&BaseUrlConf)

    baseShorterUrl := BaseUrlConf.Protocol + "://" + BaseUrlConf.Url + "/"

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        helpers.ReturnERR(w, "Failed to decode json", err)
        return
    }
    if body.Url == "" {
        helpers.ReturnERR(w, "URL can't be empty", nil)
        return
    }

    // Try to get short URL value from Redis - return if it exists
    ok, val := db.GetFromRedis(body.Url)
    if !ok {
        helpers.ReturnERR(w, val, nil)
        return
    }

    // Value is not on redis - try to fetch from DB
    if val == "" {
        ok, val = db.GetFromPostgres("short", "url_to_short", "real_url", body.Url)
        if !ok {
            helpers.ReturnERR(w, val, nil)
            return
        }
    }
    if val != "" {
        log.Println("Url already exists, returning its value!")
        helpers.ReturnOK(w, baseShorterUrl + val)
    } else {
        // Generate short URL
        ok, shorterValue := helpers.ShortenUrl(body.Url)
        if !ok {
            helpers.ReturnERR(w, "Error while trying to hash", nil)
            return
        }

        // Insert result to DB
        ok = db.InsertToRedis(body.Url, shorterValue)
        ok2 := db.InsertToRedis(shorterValue, body.Url)
        if !ok || !ok2 {
            helpers.ReturnERR(w, "Error while trying to insert value to redis DB", nil)
            return
        }

        ok = db.InsertToPostgres("url_to_short", [2]string{"real_url", "short"}, [2]string{body.Url, shorterValue})
        ok2 = db.InsertToPostgres("short_to_url", [2]string{"short", "real_url"}, [2]string{shorterValue, body.Url})
        if !ok {
            helpers.ReturnERR(w, "Error while trying to insert value to postgres DB", nil)
        } else {
            helpers.ReturnOK(w, baseShorterUrl + shorterValue)
        }
    }
    
}
