package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"twidder/router"
)

func TestPing(t *testing.T) {
	r := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected ping status 200 but got %v", w.Code)
	}
}
