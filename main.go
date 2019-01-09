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

// Todo: Fetch dynamically
var domain = "https://redirdev.herokuapp.com"

var salt = "To be fetched from a better place"
var minLength = 20

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", postHandler).Methods("POST")
	r.HandleFunc("/{token}", redirectHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"

	}
	InitDb(os.Getenv("DATABASE_URL"))

	log.Println("Server started")
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

	response := map[string]string{"response": domain + "/" + token}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func errorMessage(w http.ResponseWriter, r *http.Request, status int, e error) {
	response := map[string]string{"error": e.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
