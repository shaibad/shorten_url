package handlers

import (
    "log"
    "time"
    "encoding/json"
    "net/http"

    "url-shortener/config"
    "url-shortener/helpers"
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
    val, err := redisClient.Get(body.Url)
    if err != nil {
        helpers.ReturnERR(w, val, nil)
        return
    }

    // Value is not on redis - try to fetch from DB
    if val == "" {
        val, err = dbClient.FindByPkey("short", "url_to_short", "real_url", body.Url)
        if err != nil {
            helpers.ReturnERR(w, val, err)
            return
        }
    }
    if val != "" {
        log.Println("Url already exists, returning its value!")
        helpers.ReturnOK(w, baseShorterUrl + val)
    } else {
        // Generate short URL
        var shorterValue string
        shorterValue, err = helpers.ShortenUrl(body.Url)
        if err != nil {
            helpers.ReturnERR(w, "Error while trying to hash", err)
            return
        }

        // Insert result to DB
        err = redisClient.Set(body.Url, shorterValue, 1 * time.Hour)
        err2 := redisClient.Set(shorterValue, body.Url, 1 * time.Hour)
        if err2 != nil || err != nil {
            helpers.ReturnERR(w, "Error while trying to insert value to redis DB", err)
            return
        }

        err = dbClient.Insert("url_to_short", [2]string{"real_url", "short"}, [2]string{body.Url, shorterValue})
        err2 = dbClient.Insert("short_to_url", [2]string{"short", "real_url"}, [2]string{shorterValue, body.Url})
        if err != nil || err2 != nil{
            helpers.ReturnERR(w, "Error while trying to insert value to postgres DB", err)
        } else {
            helpers.ReturnOK(w, baseShorterUrl + shorterValue)
        }
    }
}

