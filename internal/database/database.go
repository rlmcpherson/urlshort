package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

type DB interface {
	Encode(url string) (string, error)
	Decode(url string) (string, error)
}

type pgDB struct {
	*sql.DB
}

type ErrNotFound struct {
	msg string
}

func (err ErrNotFound) Error() string {
	return err.msg
}

func NewPGDB(dbinfo string) (DB, error) {
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	return &pgDB{db}, db.Ping()
}

func (db *pgDB) Encode(url string) (string, error) {
	// create candidate randomURL and try to store until unique key
	for {
		short := randomURL()
		_, err := db.Exec("INSERT INTO urls(short, url) VALUES($1,$2);", short, url)
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == "23505" {
			// unique pk violation, try again
			continue
		}
		return short, err
	}
}

func (db *pgDB) Decode(shortURL string) (string, error) {
	var url string
	err := db.QueryRow("SELECT url FROM urls WHERE short=$1", shortURL).Scan(&url)
	if err == sql.ErrNoRows {
		return "", ErrNotFound{fmt.Sprintf("no url found for %s", shortURL)}
	}
	return url, err
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const urlChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const urlLen = 7 // shortened urls have a fixed length of 7

// randomURL returns
func randomURL() string {
	b := make([]byte, urlLen)
	for i := range b {
		b[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return string(b)
}
