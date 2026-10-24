package httpx

import (
	"net/http"
	"strconv"
)

// PathID parses the positive integer {id} path value, writing a 400 response
// and returning false when it is invalid.
func PathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		BadRequest(w, "invalid id")
		return 0, false
	}
	return id, true
}

// QueryID parses an optional positive integer query parameter. It returns nil
// when the parameter is absent and writes a 400 response when it is invalid.
func QueryID(w http.ResponseWriter, r *http.Request, name string) (*int64, bool) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return nil, true
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		BadRequest(w, "invalid "+name)
		return nil, false
	}
	return &id, true
}
