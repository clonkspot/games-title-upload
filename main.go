package main

import (
	"net/http"
	"os"
	"io"
	"log"
	"fmt"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	os.MkdirAll(filepath.Dir(path), 0666)

	out, err := os.Create(path)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Unable to create the file for writing.")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	written, err := io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
	}

	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/", uploadHandler)

	log.Fatal(http.ListenAndServe(":63230", nil))
}
