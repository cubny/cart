package jsonerror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorType is used for error codes
type errorType int

var (
	// Dear reviewer, these error codes are trivial just to demonstrate that documented
	// error codes for Rest API can help users understand and troubleshoot more efficiently
	// note that there is no standard for these numbers, they're usually defined by
	// convention in the company
	errInternalError errorType = 100500
	errUnauthorised  errorType = 100401
	errBadRequest    errorType = 100400
	errInvalidParams errorType = 100422
	errNotFound      errorType = 100404
)

// JsonError is used to return http errors encoded in json
type JsonError struct {
	// Code of the error
	Code int `json:"code"`
	// Details of the error
	Details string `json:"details"`
}

func New(errorType errorType, details string) JsonError {

	e := JsonError{Code: int(errorType)}

	switch errorType {
	case errInternalError:
		e.Details = "Internal error"
	case errBadRequest:
		e.Details = "Bad Request"
	case errUnauthorised:
		e.Details = "Unauthorised access"
	case errInvalidParams:
		e.Details = "Invalid params"
	case errNotFound:
		e.Details = "Not found"
	default:
		e.Code = 100999
		e.Details = "Unknown error"
	}

	if details != "" {
		e.Details = fmt.Sprintf("%s - %s", e.Details, details)
	}

	return e
}

func (e JsonError) write(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := struct {
		Err JsonError `json:"error"`
	}{
		Err: e,
	}

	return json.NewEncoder(w).Encode(resp)
}

// InternalError writes the error details in json with the provided details
func InternalError(w http.ResponseWriter, details string) error {
	return New(errInternalError, details).write(w, http.StatusInternalServerError)
}

// Unauthorised writes the unauthorised error details in json with the provided details
func Unauthorised(w http.ResponseWriter, details string) error {
	return New(errUnauthorised, details).write(w, http.StatusUnauthorized)
}

// BadRequest writes the BadRequest error details in json with the provided details
func BadRequest(w http.ResponseWriter, details string) error {
	return New(errBadRequest, details).write(w, http.StatusBadRequest)
}

// InvalidParams writes the UnprocessableEntity error details in json with the provided details
func InvalidParams(w http.ResponseWriter, details string) error {
	return New(errInvalidParams, details).write(w, http.StatusUnprocessableEntity)
}

// NotFound writes the NotFound error details in json with the provided details
func NotFound(w http.ResponseWriter, details string) error {
	return New(errNotFound, details).write(w, http.StatusNotFound)
}
