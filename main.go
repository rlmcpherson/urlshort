// urlshort: url shortening service

// run with -help for usage

// depends on go version 1.5 or later

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rlmcpherson/urlshort/internal/database"
)

const (
	encodePath   = "/encode/"
	decodePath   = "/decode/"
	redirectPath = "/"
)

func main() {
	// load options from env

	// connect to db

	var db database.DB

	// route handlers
	http.Handle(encodePath, encodeHandler(db))
	http.Handle(decodePath, decodeHandler(db))
	http.Handle(redirectPath, redirectHandler(db))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", opts.Port), nil))
}
