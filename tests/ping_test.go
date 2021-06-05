package test

import (
	"net/http"
	"net/http/httptest"
)

func (suite *Suite) TestPing() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	suite.NoError(err)
	suite.router.ServeHTTP(w, req)
	suite.Equal(200, w.Code, "ping status should be 200")
}
