package handler

import (
	"encoding/json"
	"net/http"
)

// EncodeJSONError .
func EncodeJSONError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// custom headers
	if headerer, ok := err.(interface{ Headers() http.Header }); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	// custome code
	if sc, ok := err.(interface{ StatusCode() int }); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	// enforce json response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
