package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", uploadHandler)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
