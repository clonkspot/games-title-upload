/* Title.png downloads */

package main

import (
	"fmt"
	"net/http"
	"os"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	if fileExists(path) {
		servePath(w, path)
	} else if fileExists(uploadPrefix + path) {
		servePath(w, uploadPrefix+path)
	} else {
		w.WriteHeader(404)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	fmt.Println(path, err)
	return err == nil
}

// Instructs nginx to serve the given path.
func servePath(w http.ResponseWriter, path string) {
	w.Header().Add("X-Accel-Redirect", internalPrefix+path)
	w.WriteHeader(200)
}
