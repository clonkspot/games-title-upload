package main

import (
	"log"
	"net/http"
	"os"
)

// The prefix which will be used in X-Accel-Redirect.
var internalPrefix = defaultValue(os.Getenv("INTERNAL_PREFIX"), "/internal/")
// The upload directory.
var uploadPrefix = defaultValue(os.Getenv("UPLOAD_PREFIX"), "./incoming")

func defaultValue(val, def string) string {
	if val == "" {
		return def
	} else {
		return val
	}
}

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
