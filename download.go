/* Title.png downloads */

package main

import (
	"net/http"
	"net/url"
	"os"
)

// Cache scenario hashes for one week.
const memoizeExpiry = 60 * 60 * 24 * 7

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	if fileExists(path) {
		servePath(w, path)
		memoizePath(r.URL, path)
	} else if p := uploadPrefix + path; fileExists(p) {
		servePath(w, p)
		memoizePath(r.URL, p)
	} else if p, ok := retrievePath(r.URL); ok {
		servePath(w, p)
	} else {
		w.WriteHeader(404)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Instructs nginx to serve the given path.
func servePath(w http.ResponseWriter, path string) {
	w.Header().Add("X-Accel-Redirect", internalPrefix+path)
	w.WriteHeader(200)
}

// Memoize the given path for the url via Redis.
func memoizePath(url *url.URL, path string) {
	client, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.Put(client)

	query := url.Query()
	if hash := query.Get("hash"); hash != "" {
		client.Cmd("SETEX", "games-title:hash:"+hash, memoizeExpiry, path)
	}
}

// Tries to retrieve the path via Redis.
func retrievePath(url *url.URL) (path string, ok bool) {
	client, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.Put(client)

	query := url.Query()
	if hash := query.Get("hash"); hash != "" {
		path, err = client.Cmd("GET", "games-title:hash:"+hash).Str()
		if err == nil && path != "" {
			ok = true
		}
	}
	return
}
