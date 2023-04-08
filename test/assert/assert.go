package assert

import (
	"reflect"
	"testing"
)

func HasLength[T any](t *testing.T, list []T, expectedLength int) {
	if length := len(list); length != expectedLength {
		t.Errorf("expected list to be of length %d, but got %d", expectedLength, length)
	}
}

func IsNil(t *testing.T, expectedValue any) {
	if expectedValue != nil {
		t.Errorf("expected nil, but got \"%v\"", expectedValue)
	}
}

func Equals[T comparable](t *testing.T, value T, expectedValue T) {
	if value != expectedValue {
		t.Errorf("expected \"%v\", but got \"%v\"", expectedValue, value)
	}
}

func DeepEquals(t *testing.T, value any, expectedValue any) {
	if !reflect.DeepEqual(value, expectedValue) {
		reflectType := reflect.TypeOf(expectedValue)
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}
		if typeName := reflectType.Name(); typeName != "" {
			t.Errorf("expected %s %v, but got %v", typeName, expectedValue, value)
			return
		}
		t.Errorf("expected %v, but got %v", expectedValue, value)
	}
}

func ErrorHasMessage(t *testing.T, err error, expectedMessage string) {
	if err == nil {
		t.Error("expected error, but got nil")
		return
	}
	if err.Error() != expectedMessage {
		t.Errorf("expected error %v, but got %v", err.Error(), expectedMessage)
	}
}
