package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type request struct {
	path string
}

// Todo: Fetch dynamically
var domain = "shit.ly"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", postHandler).Methods("POST")
	r.HandleFunc("/{token}", redirectHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	http.Redirect(w, r, fmt.Sprintf("https://www.google.com/search?q=%s", vars["token"]), 301)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var out request
	err := json.NewDecoder(r.Body).Decode(&out)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var alias = createAliasLink(out.path)

	response := map[string]string{"shortened": domain + "/" + alias}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func createAliasLink(longPath string) string {
	return "32m23"
}
