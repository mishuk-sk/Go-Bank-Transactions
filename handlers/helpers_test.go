package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHttpHandler struct {
	w http.ResponseWriter
	r *http.Request
}

func (h testHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.Write(body)
	return
}
func TestVerboseFunc(t *testing.T) {
	loggingFunc := verboseMiddleware(testHttpHandler{})
	switch loggingFunc.(type) {
	case http.HandlerFunc:
		break
	default:
		t.Errorf("verboseMiddlevare should return http.HandlerFunc not %T", loggingFunc)
	}
	body := `{
		"some_data"
	}`
	r := httptest.NewRequest(http.MethodGet, "/someurl/", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	loggingFunc.ServeHTTP(w, r)
	if w.Body.String() != body {
		t.Errorf("Can't read request body after logging func execution")
	}
}
