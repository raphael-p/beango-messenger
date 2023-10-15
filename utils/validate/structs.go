package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
)

func traverseStructFields(
	reflectValue reflect.Value,
	jsonPath string,
	missingFields []string,
) []string {
	isJSONField := regexp.MustCompile(`^JSONField\[.+\]$`)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
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
			isSet := jsonField.FieldByName("IsSet").Bool()
			isNull := jsonField.FieldByName("IsNull").Bool()
			value := jsonField.FieldByName("Value")
			isZero := value.IsZero()
			isStruct := value.Type().Kind() == reflect.Struct
			if (!isSet && !isOptional) ||
				(isNull && !isNullable) ||
				(isZero && !isZeroable && !isStruct) {
				missingFields = append(missingFields, jsonPath+jsonName)
			}
			if isStruct {
				missingFields = traverseStructFields(
					value,
					jsonPath+jsonName+".",
					missingFields,
				)
			}
		} else if field.Type.Kind() == reflect.Struct {
			missingFields = traverseStructFields(
				reflectValue.FieldByName(field.Name),
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

// Finds any fields from the struct which are zero-valued.
// `value` must be a `struct` or its pointer.
// Accepted tags:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> for JSONField only, allows field to not be set
//	`nullable:"true"` -> for JSONField only, allows field to be set to null
//	`zeroable:"true"` -> for JSONField only, allows field to be set to zero-value
//
// Used to validate the deserialisation of a JSON document.
func StructFromJSON(value any) ([]string, error) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected `value` to be a struct or its pointer, got %T", value)
	}
	return traverseStructFields(reflectValue, "", []string{}), nil
}

// Type wrapper used to validate the struct which a JSON document is
// unmarshaled to. Accepted tags:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> allows field to not be set
//	`nullable:"true"` -> allows field to be set to null
//	`zeroable:"true"` -> allows field to be set to zero-value
type JSONField[T any] struct {
	Value  T
	IsNull bool
	IsSet  bool
}

// this is called implicitly when unmarshalling into a struct containing JSONField
func (i *JSONField[any]) UnmarshalJSON(data []byte) error {
	i.IsSet = true

	if string(data) == "null" {
		i.IsNull = true
		return nil
	}

	var temp any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.IsNull = false
	return nil
}

// Validates that the input is a pointer to a struct where all fields are strings.
func PointerToStringStruct(ptr any) bool {
	reflectValue := reflect.ValueOf(ptr)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return false
	}

	reflectType := reflectValue.Type().Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		if reflectType.Field(i).Type.Kind() != reflect.String {
			return false
		}
	}

	return true
}
