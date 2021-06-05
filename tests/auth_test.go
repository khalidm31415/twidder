package test

import (
	"bytes"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyStringLongerThan struct {
	n int
}

func (a AnyStringLongerThan) Match(v driver.Value) bool {
	s, ok := v.(string)
	if ok && len(s) >= 50 {
		return true
	}
	return false
}

func (suite *Suite) TestSignup() {
	req, err := http.NewRequest("POST", "/auth/signup", nil)
	suite.NoError(err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(400, w.Code, "signup status should be 400 when given empty data")

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs("test", "Test User", "test@test.com", AnyStringLongerThan{50}, AnyTime{}, AnyTime{}, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	jsonStr := []byte(`{"email": "test@test.com", "userName": "test", "password": "TestPassword", "displayName": "Test User"}`)
	req, err = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonStr))
	suite.NoError(err)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(200, w.Code, "signup status should be 200 when given valid data")

	jsonStr = []byte(`{"email": "test", "userName": "tes", "password": "test", "displayName": "test"}`)
	req, err = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonStr))
	suite.NoError(err)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(400, w.Code, "signup status should be 400 when given invalid data (email, userName, and password too short)")
}

func (suite *Suite) TestLogin() {
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("SELECT FROM `users`")).
		WithArgs("test", AnyStringLongerThan{50}).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	jsonStr := []byte(`{"userName": "test", "password": "TestPassword"}`)
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	suite.NoError(err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	jsonStr = []byte(`{"userName": "tes"}`)
	req, err = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	suite.NoError(err)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(401, w.Code, "signin status should be 401 when password is not given")
}

func (suite *Suite) TestLogout() {
	req, err := http.NewRequest("POST", "/auth/logout", nil)
	suite.NoError(err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(200, w.Code, "signout status should be 200")
}
