package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

type request struct {
	Path string
}

var prefix string
var minLength = 5

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", postHandler).Methods("POST")
	r.HandleFunc("/{token}", redirectHandler).Methods("GET")

	port := getEnv("PORT", "8888")
	prefix = getEnv("PATH_PREFIX", "localhost:"+port)
	// Note: Database_url is the one configuration you need to set before trying this out.
	InitDb(getEnv("DATABASE_URL", "postgres://user:pw@host:5432/nameofdb"))
	InitHash(getEnv("HASH_SALT", "Unsafe salt"))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(":"+port, r))

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, hasherr := Dehash(vars["token"])
	if hasherr != nil {
		errorMessage(w, r, http.StatusBadRequest, hasherr)
		return
	}

	longURL, err := RetrieveLongURL(id)
	if err != nil {
		errorMessage(w, r, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, longURL, 301)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var req request

	jsonerr := json.NewDecoder(r.Body).Decode(&req)
	if jsonerr != nil {
		log.Println(jsonerr)
		errorMessage(w, r, http.StatusBadRequest, errors.New("Bad request"))
		return
	}

	_, urlErr := url.ParseRequestURI(req.Path)
	if urlErr != nil {
		errorMessage(w, r, http.StatusBadRequest, errors.New("Bad url. Forgot http://?"))
		return
	}

	token, err := SubmitLongURL(req.Path)
	if err != nil {
		errorMessage(w, r, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"response": prefix + "/" + token}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func errorMessage(w http.ResponseWriter, r *http.Request, status int, e error) {
	response := map[string]string{"error": e.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
