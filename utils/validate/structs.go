package validate

import (
	"reflect"
)

// Finds any fields from the struct's type that are not in the struct itself.
// `ptr` must point to a `struct` where fields may have `json` and `optional` tags.
// Used to validate the deserialisation of a JSON document.
func PointerToStructFromJSON(ptr any) []string {
	reflectValue := reflect.ValueOf(ptr).Elem()
	reflectType := reflectValue.Type()
	return traverseStructFields(reflectValue, reflectType, "", []string{})
}

func traverseStructFields(
	reflectValue reflect.Value,
	reflectType reflect.Type,
	jsonPath string,
	missingFields []string,
) []string {
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		if optional := field.Tag.Get("optional"); optional == "true" {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		if field.Type.Kind() == reflect.Struct {
			newJsonPath := jsonPath + tag + "."
			missingFields = traverseStructFields(
				reflectValue.FieldByName(field.Name),
				field.Type,
				newJsonPath,
				missingFields,
			)
		} else {
			if reflectValue.FieldByName(field.Name).IsZero() {
				fullFieldPath := jsonPath + tag
				missingFields = append(missingFields, fullFieldPath)
			}
		}
	}
	return missingFields
}

// Validates that the input is a pointer to a struct
// where all fields are strings.
func PointerToStringStruct(ptr any) bool {
	reflectValue := reflect.ValueOf(ptr)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return false
	}

	reflectType := reflectValue.Elem().Type()
	for i := 0; i < reflectType.NumField(); i++ {
		if reflectType.Field(i).Type.Kind() != reflect.String {
			return false
		}
	}

	return true
}
