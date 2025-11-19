package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

// RespondWithJSON sends a JSON response with the specified HTTP status code.
// It accepts an optional payload that will be encoded as JSON in the response body.
// If encoding fails, it returns a 500 Internal Server Error.
func RespondWithJSON(w http.ResponseWriter, status int, payload ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(payload) > 0 {
		err := json.NewEncoder(w).Encode(payload[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// RespondWithError sends a JSON error response with the specified HTTP status code.
// It accepts an optional custom error message, defaulting to a generic message if none provided.
// The response format is {"message": "error text"}.
func RespondWithError(w http.ResponseWriter, code int, msg ...string) {
	if len(msg) > 0 {
		RespondWithJSON(w, code, map[string]any{"message": msg[0]})
	} else {
		RespondWithJSON(w, code, map[string]any{"message": "unexpected error, try again!"})
	}
}

// ParseBody decodes the JSON request body into the provided interface.
// It automatically closes the request body after reading.
// Returns an error if the JSON is malformed or cannot be decoded.
func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetScanFields returns a slice of pointers to struct fields for scanning SQL results.
// It uses reflection to iterate through all struct fields and return their addresses.
// Pass a pointer to the struct.
// Example: GetScanFields(&user) => []interface{}{&user.ID, &user.FirstName, &user.LastName, ...}
func GetScanFields(s interface{}) []interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		log.Fatal("Input must be a pointer to a struct")
	}
	val = val.Elem() // like &user => user

	fields := make([]interface{}, val.NumField()) // like user => user.ID, user.FirstName, ...
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}
	return fields
}

// GetExecFields returns a slice of struct field values, excluding specified fields.
// It uses reflection to extract field values, useful for SQL INSERT/UPDATE operations.
// Excluded fields are typically auto-generated (ID, timestamps) or should not be modified.
// Example: GetExecFields(user, "ID", "CreatedAt") => []interface{}{user.FirstName, user.LastName, ...}
func GetExecFields(s interface{}, excludeFields ...string) []interface{} {
	val := reflect.ValueOf(s) // like user => user.FirstName, user.LastName, ...
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		log.Fatal("Input must be a struct or pointer to a struct")
	}

	// Convert excluded fields to a map for fast lookup
	excluded := make(map[string]bool)
	for _, field := range excludeFields {
		excluded[field] = true
	}

	var fields []interface{}
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if !excluded[fieldName] {
			fields = append(fields, val.Field(i).Interface())
		}
	}
	return fields
}
