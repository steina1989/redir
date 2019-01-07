package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type request struct {
	path string
}

// Todo: Fetch dynamically
var domain = "redirdev.herokuapp.com"

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
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL, _ := retrieveLongURL(vars["token"])
	http.Redirect(w, r, longURL, 301)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var alias, _ = submitLongURL(req.path)

	response := map[string]string{"shortened": domain + "/" + alias}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
