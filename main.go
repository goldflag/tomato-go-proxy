package main

import (
	"net/http"

	"tomato_proxy/endpoints"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.HandleFunc("/fetchPlayer/{server}/{id}", endpoints.FetchPlayer)

	http.ListenAndServe(":8000", r)
}
