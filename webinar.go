package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "user-agent %s", r.Header.Get("user-agent"))
	})
	http.Handle("/", r)
	log.Println("starting server...")
	log.Println(http.ListenAndServe(":8000", r))
}
