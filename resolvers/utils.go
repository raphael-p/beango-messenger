package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/raphael-p/beango/utils"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer
func bindRequestJSON(w *utils.ResponseWriter, r *http.Request, ptr any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		response := fmt.Sprintf("malformed request body: %v", err)
		w.StringResponse(http.StatusBadRequest, response)
		return false
	}
	if fields := findMissingFields(ptr); len(fields) != 0 {
		response := fmt.Sprintf("missing required field(s): %s", fields)
		w.StringResponse(http.StatusBadRequest, response)
		return false

	}
	return true
}

type requiredField struct {
	name string
	json string
}

func (rf requiredField) String() string {
	if rf.json != "" {
		return rf.json
	}
	return rf.name
}

// finds fields which are defined in a struct's type but are not in the struct
// and do not have the optional tag
func findMissingFields(ptr any) []string {
	// fetch required fields (fields are required by default)
	reflectValue := reflect.ValueOf(ptr).Elem()
	reflectType := reflectValue.Type()
	requiredFields := make([]requiredField, 0, reflectType.NumField())
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		if optional := field.Tag.Get("optional"); optional == "true" {
			continue
		}
		requiredField := requiredField{field.Name, field.Tag.Get("json")}
		requiredFields = append(requiredFields, requiredField)
	}

	// look for required fields missing from the struct
	missingFields := make([]string, 0, len(requiredFields))
	for _, field := range requiredFields {
		if reflectValue.FieldByName(field.name).IsZero() {
			missingFields = append(missingFields, field.String())
		}
	}
	return missingFields
}
