/* Title.png uploads */

package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const filesizeLimit = 100 * 1024

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check that the body is a PNG file with a maximum of 100KB.
	var buf bytes.Buffer

	n, err := buf.ReadFrom(io.LimitReader(r.Body, filesizeLimit))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	if n >= filesizeLimit {
		w.WriteHeader(400)
		fmt.Fprintln(w, "File too large.")
		return
	}

	_, err = png.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "File not a PNG image.")
		return
	}

	path := "." + r.URL.Path
	os.MkdirAll(filepath.Dir(path), 0777)

	out, err := os.Create(path)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Unable to create the file for writing.")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, &buf)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	w.WriteHeader(200)
	log.Println("->", path, "\t", n)
}
