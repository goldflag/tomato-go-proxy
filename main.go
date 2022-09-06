package main

import (
	"fmt"
	"net/http"

	"tomato_proxy/endpoints"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")

}

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.HandleFunc("/fetchPlayer/{server}/{id}", endpoints.FetchPlayer)

	// fmt.Println(os.Getenv("FOO"))
	// fmt.Println(endpoints.FetchPlayer())

	// http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8000", r)
}
