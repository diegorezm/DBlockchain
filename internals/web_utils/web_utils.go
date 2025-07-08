package webutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// --- Client-side utilities ---

// ParseJSON is a generic function to parse JSON from an io.ReadCloser into a target type.
// It automatically closes the provided ReadCloser.
// T is the type parameter, representing the expected structure of the JSON data.
func ParseJSON[T any](body io.ReadCloser) (T, error) {
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("http_utils: failed to read response body: %w", err)
	}

	var result T

	// Unmarshal JSON bytes into the address of the result variable.
	// json.Unmarshal requires a pointer. If T is a value type (e.g., `[]Block`),
	// `&result` gives its address. If T is already a pointer type (e.g., `*MyStruct`),
	// `&result` would be a pointer to a pointer, but json.Unmarshal handles this correctly
	// by dereferencing the pointer type itself.
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("http_utils: failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

// --- Server-side utilities ---

type JSONResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// WriteJSON writes a generic data structure as a JSON response to the http.ResponseWriter.
// It sets the Content-Type header to application/json and writes the provided HTTP status code.
func WriteJSON[T any](w http.ResponseWriter, statusCode int, data T, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := JSONResponse[T]{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("http_utils: failed to encode JSON response: %v", err), http.StatusInternalServerError)
	}
}

// WriteError writes a structured error response as JSON to the http.ResponseWriter.
// It sets the Content-Type header to application/json and writes the provided HTTP status code.
func WriteError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := JSONResponse[any]{
		Status: "error",
		Error:  errMsg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("http_utils: failed to encode error JSON response: %v", err), http.StatusInternalServerError)
	}
}

// --- Convenience functions for common HTTP responses ---

func WriteSuccess[T any](w http.ResponseWriter, data T, message string) {
	WriteJSON(w, http.StatusOK, data, message)
}

func WriteBadRequest(w http.ResponseWriter, errMsg string) {
	WriteError(w, http.StatusBadRequest, errMsg)
}

func WriteUnauthorized(w http.ResponseWriter, errMsg string) {
	WriteError(w, http.StatusUnauthorized, errMsg)
}

func WriteNotFound(w http.ResponseWriter, errMsg string) {
	WriteError(w, http.StatusNotFound, errMsg)
}

func WriteInternalServerError(w http.ResponseWriter, errMsg string) {
	WriteError(w, http.StatusInternalServerError, errMsg)
}

func WriteCustomError(w http.ResponseWriter, statusCode int, errMsg string) {
	WriteError(w, statusCode, errMsg)
}
