package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	expTime   = 60
	longURL   = "https://www.baidu.com"
	shortLink = "IFHzaO"
)

type storageMock struct {
	mock.Mock
}

var app App
var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) UnShorten(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortLinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func init() {
	app = App{}
	mockR = new(storageMock)
	app.Init(&Env{S: mockR})
}

func TestCreateShortLink(t *testing.T) {
	// define request body content
	var jsonContent = []byte(`{
		"url": "https://www.baidu.com",
		"expire": 60
	}`)

	// create a request
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonContent))
	if err != nil {
		t.Fatalf("could not create the request, error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// setup expectations
	mockR.On("Shorten", longURL, expTime).Return(shortLink, nil).Once()
	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK, "don't get the expected status code")

	resp := struct {
		ShortLink string
	}{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode the response, error: %v", err)
	}

	assert.Equal(t, resp.ShortLink, shortLink, "don't get the expected short link")
}

func TestRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortLink)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Fatalf("could not create request, error: %v", err)
	}

	mockR.On("UnShorten", shortLink).Return(longURL, nil).Once()
	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	assert.Equal(t, rw.Code, http.StatusTemporaryRedirect, "don't get the expected status code")
}
