package main

import (
	"log"
	"net/http"
	"os"
	redis "github.com/fzzy/radix/extra/pool"
)

// The prefix which will be used in X-Accel-Redirect.
var internalPrefix = defaultValue(os.Getenv("INTERNAL_PREFIX"), "/internal/")

// The upload directory.
var uploadPrefix = defaultValue(os.Getenv("UPLOAD_PREFIX"), "./incoming/")

var redisNetwork = defaultValue(os.Getenv("REDIS_NETWORK"), "tcp")
var redisAddress = defaultValue(os.Getenv("REDIS_ADDRESS"), "127.0.0.1:6379")
var redisPool *redis.Pool

func defaultValue(val, def string) string {
	if val == "" {
		return def
	} else {
		return val
	}
}

func initRedis() {
	var err error
	redisPool, err = redis.NewPool(redisNetwork, redisAddress, 5)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initRedis()

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
	log.Print("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
