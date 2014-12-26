package main

import (
	"log"
	"net/http"
	"os"
)

// The prefix which will be used in X-Accel-Redirect.
var internalPrefix = "/serve/"
// The upload directory.
var uploadPrefix = "./incoming/"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			downloadHandler(w, r)
		case "POST":
			uploadHandler(w, r)
		default:
			w.WriteHeader(400)
		}
	})

	port := os.Getenv("PORT")
	log.Print("Listening on port "+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
