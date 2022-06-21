package handlers

import (
    "net/http"
    "github.com/gorilla/mux"

    "url-shortener/helpers"
)

func GetUrlHandler(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
	short := vars["short"]

    if short == "" {
        helpers.ReturnERR(w, "You must provide the short path", nil)
        return
    }

    // Try to get short URL value from Redis - return if it exists
    val, err := redisClient.Get(short)
    if err != nil {
        helpers.ReturnERR(w, val, nil)
        return
    }

    // Value is not on redis - try to fetch from DB
    if val == "" {
        val, err = dbClient.FindByPkey("real_url", "short_to_url", "short", short)
        if err != nil {
            helpers.ReturnERR(w, val, err)
            return
        }
    }

	if val != "" {
        http.Redirect(w, r, val, http.StatusSeeOther)
    } else {
        helpers.ReturnERR(w, "Url not found", nil)
    }
}
