package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRoute(t *testing.T) {
	router := setupRouter(&apiConfig{})               // health doesnt touch the store, so no DB needed
	req := httptest.NewRequest("GET", "/health", nil) //a fake GET /health request
	rec := httptest.NewRecorder()                     // a fake resonse, records whats written
	router.ServeHTTP(rec, req)                        // run the requet thru the router
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	
}
func TestUnknownRoute(t *testing.T) {
	router := setupRouter(&apiConfig{})
	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
