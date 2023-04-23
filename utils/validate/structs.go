package validate

import (
	"encoding/json"
	"reflect"
	"regexp"
)

// Finds any fields from the struct's type that are not in the struct itself.
// `ptr` must point to a `struct`.
// Accepted tags:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> for JSONField only, allows field to not be set
//	`nullable:"true"` -> for JSONField only, allows field to be set to null
//	`zeroable:"true"` -> for JSONField only, allows field to be set to zero-value
//
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
	isJSONField := regexp.MustCompile(`^JSONField\[.+\]$`)
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tags := field.Tag

		jsonName := tags.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		}

		if isJSONField.MatchString(field.Type.Name()) {
			jsonField := reflectValue.FieldByName(field.Name)
			isOptional := tags.Get("optional") == "true"
			isNullable := tags.Get("nullable") == "true"
			isZeroable := isOptional || isNullable || tags.Get("zeroable") == "true"
			isSet := jsonField.FieldByName("Set").Bool()
			isNull := jsonField.FieldByName("Null").Bool()
			value := jsonField.FieldByName("Value")
			_type := value.Type()
			isZero := value.IsZero()
			isStruct := _type.Kind() == reflect.Struct
			if (!isSet && !isOptional) ||
				(isNull && !isNullable) ||
				(isZero && !isZeroable && !isStruct) {
				missingFields = append(missingFields, jsonPath+jsonName)
			}
			if isStruct {
				missingFields = traverseStructFields(
					value,
					_type,
					jsonPath+jsonName+".",
					missingFields,
				)
			}
		} else if field.Type.Kind() == reflect.Struct {
			missingFields = traverseStructFields(
				reflectValue.FieldByName(field.Name),
				field.Type,
				jsonPath+jsonName+".",
				missingFields,
			)
		} else {
			if reflectValue.FieldByName(field.Name).IsZero() {
				missingFields = append(missingFields, jsonPath+jsonName)
			}
		}
	}
	return missingFields
}

// Type wrapper used to validate the struct which a JSON document is
// unmarshaled to. Accepted tags:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> allows field to not be set
//	`nullable:"true"` -> allows field to be set to null
//	`zeroable:"true"` -> allows field to be set to zero-value
type JSONField[T any] struct {
	Value T
	Null  bool
	Set   bool
}

// this is called implicitly when unmarshalling into a struct containing JSONField
func (i *JSONField[any]) UnmarshalJSON(data []byte) error {
	i.Set = true

	if string(data) == "null" {
		i.Null = true
		return nil
	}

	var temp any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Null = false
	return nil
}

// Validates that the input is a pointer to a struct where all fields are strings.
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
