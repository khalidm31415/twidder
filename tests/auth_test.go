package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"twidder/router"
)

func TestSingup(t *testing.T) {
	r := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/signup", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected signup status 400 with empty data but got %v", w.Code)
	}
}
