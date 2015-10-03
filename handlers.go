package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rlmcpherson/urlshort/internal/database"
)

func decodeHandler(db database.DB) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "use GET to decode urls", http.StatusMethodNotAllowed)
			return
		}

		shorturl := r.URL.Path[len(decodePath):]
		url, err := db.Decode(shorturl)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeLogErr(w, []byte(fmt.Sprintf("url not found: %s", err)))
			return
		}
		writeLogErr(w, []byte(url))
	}

	return http.HandlerFunc(handler)
}

func encodeHandler(db database.DB) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeLogErr(w, []byte("use POST to encode urls"))
			return
		}
		if r.Body == nil {
			http.Error(w, http.ErrShortBody.ErrorString, http.StatusBadRequest)
			return
		}

		longurl, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shorturl, err := db.Encode(string(longurl))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeLogErr(w, []byte(fmt.Sprintf("error encoding %s: %s", longurl, err)))
			return
		}
		writeLogErr(w, []byte(shorturl))
	}

	return http.HandlerFunc(handler)

}

func redirectHandler(db database.DB) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeLogErr(w, []byte("use GET to decode urls"))
		}

		shorturl := r.URL.Path[len(redirectPath):]
		url, err := db.Decode(shorturl)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeLogErr(w, []byte(fmt.Sprintf("url not found: %s", err)))
		}
		http.Redirect(w, r, string(url), 301)
	}

	return http.HandlerFunc(handler)
}

func writeLogErr(w io.Writer, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		log.Printf("write error: %s", err)
	}

}
