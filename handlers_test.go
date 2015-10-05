package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rlmcpherson/urlshort/internal/database"
)

// TODO: consider DRYing out test cases with func for test case loop

type fakedb struct {
	returnErr bool
}

func (db fakedb) Decode(url string) (string, error) {
	if db.returnErr {
		return "", fmt.Errorf("error decoding %s", url)
	}
	return fmt.Sprintf("decode %s", url), nil
}

func (db fakedb) Encode(url string) (string, error) {
	if db.returnErr {
		return "", fmt.Errorf("error encoding %s", url)
	}
	return fmt.Sprintf("encode %s", url), nil
}

type handlerTest struct {
	db        database.DB
	reqMethod string
	reqURL    string
	reqBody   io.Reader
	respCode  int
	respStr   string
}

func TestDecodeHandler(t *testing.T) {

	var tests = []handlerTest{
		// valid request
		{db: fakedb{returnErr: false},
			reqMethod: "GET",
			reqURL:    fmt.Sprintf("%sfoo.com", decodePath),
			reqBody:   nil,
			respCode:  http.StatusOK,
			respStr:   "decode foo.com",
		},
		// invalid method
		{db: fakedb{returnErr: false},
			reqMethod: "POST",
			reqURL:    fmt.Sprintf("%sfoo.com", decodePath),
			reqBody:   nil,
			respCode:  http.StatusMethodNotAllowed,
			respStr:   "use GET",
		},
		// storage error
		{db: fakedb{returnErr: true},
			reqMethod: "GET",
			reqURL:    fmt.Sprintf("%sfoo.com", decodePath),
			reqBody:   nil,
			respCode:  http.StatusInternalServerError,
			respStr:   "error decoding",
		},
	}

	for _, tt := range tests {
		dh := decodeHandler(tt.db)
		req, err := http.NewRequest(tt.reqMethod, tt.reqURL, tt.reqBody)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		dh.ServeHTTP(w, req)
		t.Logf("%d - %s", w.Code, w.Body.String())
		if w.Code != tt.respCode {
			t.Errorf("expected %d, got %d", tt.respCode, w.Code)
		}
		if !strings.Contains(w.Body.String(), tt.respStr) {
			t.Errorf("response: expected %s, got %s", tt.respStr, w.Body.String())
		}
	}
}

func TestEncodeHandler(t *testing.T) {

	var tests = []handlerTest{
		// valid request
		{db: fakedb{returnErr: false},
			reqMethod: "POST",
			reqURL:    encodePath,
			reqBody:   bytes.NewBufferString("foo.com"),
			respCode:  http.StatusOK,
			respStr:   "encode foo.com",
		},
		// invalid method
		{db: fakedb{returnErr: false},
			reqMethod: "GET",
			reqURL:    encodePath,
			reqBody:   nil,
			respCode:  http.StatusMethodNotAllowed,
			respStr:   "use POST",
		},
		// missing body
		{db: fakedb{returnErr: false},
			reqMethod: "POST",
			reqURL:    encodePath,
			reqBody:   nil,
			respCode:  http.StatusBadRequest,
			respStr:   http.ErrShortBody.ErrorString,
		},
		// storage error
		{db: fakedb{returnErr: true},
			reqMethod: "POST",
			reqURL:    encodePath,
			reqBody:   bytes.NewBufferString("foo.com"),
			respCode:  http.StatusInternalServerError,
			respStr:   "error encoding foo.com",
		},
	}

	for _, tt := range tests {
		dh := encodeHandler(tt.db)
		req, err := http.NewRequest(tt.reqMethod, tt.reqURL, tt.reqBody)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		dh.ServeHTTP(w, req)
		t.Logf("%d - %s", w.Code, w.Body.String())
		if w.Code != tt.respCode {
			t.Errorf("expected %d, got %d", tt.respCode, w.Code)
		}
		if !strings.Contains(w.Body.String(), tt.respStr) {
			t.Errorf("response: expected %s, got %s", tt.respStr, w.Body.String())
		}
	}
}

func TestRedirectHandler(t *testing.T) {

	var tests = []handlerTest{
		// valid request
		{db: fakedb{returnErr: false},
			reqMethod: "GET",
			reqURL:    fmt.Sprintf("%sfoo.com", redirectPath),
			reqBody:   bytes.NewBufferString("foo.com"),
			respCode:  http.StatusMovedPermanently,
			respStr:   `<a href="/decode foo.com">Moved Permanently</a>`,
		},
		// invalid method
		{db: fakedb{returnErr: false},
			reqMethod: "POST",
			reqBody:   nil,
			reqURL:    fmt.Sprintf("%sfoo.com", redirectPath),
			respCode:  http.StatusMethodNotAllowed,
			respStr:   "GET only",
		},
		// storage error
		{db: fakedb{returnErr: true},
			reqMethod: "GET",
			reqURL:    fmt.Sprintf("%sfoo.com", redirectPath),
			reqBody:   bytes.NewBufferString("foo.com"),
			respCode:  http.StatusInternalServerError,
			respStr:   "error decoding foo.com",
		},
	}

	for _, tt := range tests {
		dh := redirectHandler(tt.db)
		req, err := http.NewRequest(tt.reqMethod, tt.reqURL, tt.reqBody)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		dh.ServeHTTP(w, req)
		t.Logf("%d - %s", w.Code, w.Body.String())
		if w.Code != tt.respCode {
			t.Errorf("expected %d, got %d", tt.respCode, w.Code)
		}
		if !strings.Contains(w.Body.String(), tt.respStr) {
			t.Errorf("response: expected %s, got %s", tt.respStr, w.Body.String())
		}
	}
}
