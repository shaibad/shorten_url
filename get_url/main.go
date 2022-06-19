package main

import (
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
	"get_url/helpers"
)

func returnERR(w http.ResponseWriter, message string, err error){
    log.Println(message, err)
    m := map[string]string{
        "Status": "Error",
        "Message": message,
	}
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(m)
}

func GetUrlHandler(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
	short := vars["short"]

    if short == "" {
        returnERR(w, "You must provide the short path", nil)
        return
    }

	// Try to get short URL value from Redis - return if it exists
    ok, val := helpers.GetFromRedis(short)
    if !ok {
        returnERR(w, val, nil)
        return
    }

    // Value is not on redis - try to fetch from DB
    if val == "" {
        ok, val = helpers.GetFromPostgres("real_url", "short_to_url", "short", short)
        if !ok {
            returnERR(w, val, nil)
            return
        }
    }

	if val != "" {
        http.Redirect(w, r, val, http.StatusSeeOther)
    } else {
        returnERR(w, "Url not found", nil)
    }
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{short}", GetUrlHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", router))
}