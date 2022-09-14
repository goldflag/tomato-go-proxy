package main

import (
	"net/http"
	"os"

	"tomato_proxy/endpoints"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.HandleFunc("/fetchPlayer/{server}/{id}", endpoints.FetchPlayer)

	http.ListenAndServe(fmt.Sprint(":", os.Getenv("PORT")), r)
}
