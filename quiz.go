/* Title.png downloads */

package main

import (
	// "log"
	"net/http"
	// "net/url"
	"os"
	"path"
	"path/filepath"
	"math/rand"
	"crypto/sha256"
	crand "crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
)

type Question struct {
	Image string
	Options []string
	Answer string
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	paths, err := findScenarios()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var question Question

	question.Image = randomScenario(paths)
	answer := path.Base(question.Image)
	options := []string{answer}
	for len(options) < 4 {
		opt := path.Base(randomScenario(paths))
		if !stringInSlice(opt, options) {
			options = append(options, opt)
		}
	}
	shuffle(options)
	question.Options = options
	encrypted, err := encryptAnswer(answer)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	question.Answer = encrypted

	b, err := json.Marshal(question)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Walks the current directory to find all possible scenarios, returning their paths.
func findScenarios() ([]string, error) {
	paths := make([]string, 1)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	return paths, err
}

// Selects a random scenario from the array.
func randomScenario(scenarios []string) string {
	n := rand.Intn(len(scenarios))
	return scenarios[n]
}

const answerRandLen = 4

// "Encrypts" the answer to hand off to the client.
func encryptAnswer(answer string) (string, error) {
	buf := make([]byte, answerRandLen)
	_, err := crand.Read(buf)
	if err != nil {
		return "", err
	}
	buf = append(buf, []byte(answerSecret)...)
	buf = append(buf, []byte(answer)...)
	result := make([]byte, sha256.Size + answerRandLen)
	copy(result[:answerRandLen],  buf[:answerRandLen])
	shasum := sha256.Sum256(buf)
	copy(result[answerRandLen:], shasum[:])
	return base64.StdEncoding.EncodeToString(result), nil
}

// Verifies a previously encrypted answer.
func verifyAnswer(answer string, encrypted string) bool {
	// encbuf contains a nonce and the hash.
	encbuf, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return false
	}
	buf := make([]byte, answerRandLen)
	copy(buf, encbuf[:answerRandLen])
	buf = append(buf, []byte(answerSecret)...)
	buf = append(buf, []byte(answer)...)
	shasum := sha256.Sum256(buf)
	return subtle.ConstantTimeCompare(shasum[:], encbuf[answerRandLen:]) == 1
}

// Generics ftw
func shuffle(array []string) {
	for i := len(array); i > 1; {
		i--
		j := rand.Intn(i)
		array[i], array[j] = array[j], array[i]
	}
}

// again
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
