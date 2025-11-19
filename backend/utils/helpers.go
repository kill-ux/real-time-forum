package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

// RespondWithJSON sends a JSON response with the specified status code
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

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, code int, msg ...string) {
	if len(msg) > 0 {
		RespondWithJSON(w, code, map[string]any{"message": msg[0]})
	} else {
		RespondWithJSON(w, code, map[string]any{"message": "unexpected error, try again!"})
	}
}

// ParseBody decodes JSON request body into the provided interface
func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetScanFields returns pointers to struct fields for SQL scanning
// Example: GetScanFields(&user) => []*interface{}{&user.ID, &user.FirstName, ...}
func GetScanFields(s interface{}) []interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		log.Fatal("Input must be a pointer to a struct")
	}
	val = val.Elem()

	fields := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}
	return fields
}

// GetExecFields returns struct field values, excluding specified fields
// Example: GetExecFields(user, "ID", "CreatedAt") => []interface{}{user.FirstName, user.LastName, ...}
func GetExecFields(s interface{}, excludeFields ...string) []interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		log.Fatal("Input must be a struct or pointer to a struct")
	}

	// Create map of excluded fields for fast lookup
	excluded := make(map[string]bool)
	for _, field := range excludeFields {
		excluded[field] = true
	}

	// Collect non-excluded field values
	var fields []interface{}
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if !excluded[fieldName] {
			fields = append(fields, val.Field(i).Interface())
		}
	}
	return fields
}
