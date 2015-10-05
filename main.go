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
	if err := optsFromEnv(); err != nil {
		log.Fatal(err)
	}

	// connect to db
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		opts.DBUser, opts.DBPassword, opts.DBAddr, opts.DBName)
	db, err := database.NewPGDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	// route handlers
	http.Handle(encodePath, encodeHandler(db))
	http.Handle(decodePath, decodeHandler(db))
	http.Handle(redirectPath, redirectHandler(db))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", opts.Port), nil))
}
