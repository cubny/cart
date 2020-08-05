package jsonerror_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cubny/cart/internal/jsonerror"
	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()

	jsonerror.BadRequest(w, "test")
	assert.Equal(t, w.Code, http.StatusBadRequest)

	expectedBody := `{"error":{"code":100400, "details":"Bad Request - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()

	jsonerror.InternalError(w, "test")
	assert.Equal(t, w.Code, http.StatusInternalServerError)

	expectedBody := `{"error":{"code":100500, "details":"Internal error - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestInvalidParams(t *testing.T) {
	w := httptest.NewRecorder()

	jsonerror.InvalidParams(w, "test")
	assert.Equal(t, w.Code, http.StatusUnprocessableEntity)

	expectedBody := `{"error":{"code":100422, "details":"Invalid params - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestUnauthorised(t *testing.T) {
	w := httptest.NewRecorder()

	jsonerror.Unauthorised(w, "test")
	assert.Equal(t, w.Code, http.StatusUnauthorized)

	expectedBody := `{"error":{"code":100401, "details":"Unauthorised access - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	jsonerror.NotFound(w, "test")
	assert.Equal(t, w.Code, http.StatusNotFound)

	expectedBody := `{"error":{"code":100404, "details":"Not found - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func assertBody(t *testing.T, expectedBody string, actualBody *bytes.Buffer) {
	t.Helper()

	body, err := ioutil.ReadAll(actualBody)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedBody, string(body))
}
