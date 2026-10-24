// Package httpx contains shared HTTP request and response helpers.
package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxBodyBytes = 1 << 20

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// DecodeJSON decodes a single JSON object from the request body into dst,
// rejecting unknown fields, trailing data, and bodies larger than 1 MiB.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON object")
	}
	return nil
}

// Error writes a JSON error body of the form {"error": message}.
func Error(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// NotImplemented marks scaffold endpoints that have no feature logic yet.
func NotImplemented(w http.ResponseWriter) {
	Error(w, http.StatusNotImplemented, "not_implemented")
}

func NotLoggedIn(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "not logged in")
}

func PermissionDenied(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "permission denied")
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}
