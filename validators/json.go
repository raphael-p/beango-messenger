package validators

import (
	"reflect"
)

// Used to validate the deserialisation of a JSON document.
// Takes a pointer to a struct and finds any fields from the struct's type
// that are not in the struct itself. Excludes fields with an 'optional' tag.
func DeserialisedJSON(ptr any) []string {
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

		if field.Type.Kind() == reflect.Struct {
			newJsonPath := jsonPath + field.Tag.Get("json") + "."
			missingFields = traverseStructFields(
				reflectValue.FieldByName(field.Name),
				field.Type,
				newJsonPath,
				missingFields,
			)
		} else {
			if reflectValue.FieldByName(field.Name).IsZero() {
				fullFieldPath := jsonPath + field.Tag.Get("json")
				missingFields = append(missingFields, fullFieldPath)
			}
		}
	}
	return missingFields
}
