package assert

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func HasLength[T any](t *testing.T, list []T, expectedLength int) {
	if length := len(list); length != expectedLength {
		t.Errorf("expected list to have %d elements, got %d", expectedLength, length)
	}
}

func IsNil(t *testing.T, values ...any) {
	for _, value := range values {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			t.Errorf("expected nil, got \"%v\"", value)
		}
	}
}

func IsNotNil(t *testing.T, values ...any) {
	for _, value := range values {
		if value == nil || reflect.ValueOf(value).IsNil() {
			t.Error("expected non-nil, got nil")
		}
	}
}

func Equals[T comparable](t *testing.T, value T, expectedValue T) {
	if value != expectedValue {
		t.Errorf("expected \"%v\", got \"%v\"", expectedValue, value)
	}
}

func NotEquals[T comparable](t *testing.T, value T, expectedValue T) {
	if value == expectedValue {
		t.Errorf("expected not \"%v\"", value)
	}
}

func Contains(t *testing.T, value string, expectedValues ...string) {
	for _, expectedValue := range expectedValues {
		if !strings.Contains(value, expectedValue) {
			t.Errorf("expected to find \"%s\" in \"%s\"", expectedValue, value)
		}
	}
}

func NotContains(t *testing.T, value string, expectedValues ...string) {
	for _, expectedValue := range expectedValues {
		if strings.Contains(value, expectedValue) {
			t.Errorf("expected not to find \"%s\" in \"%s\"", expectedValue, value)
		}
	}
}

func DeepEquals(t *testing.T, value any, expectedValue any) {
	if !reflect.DeepEqual(value, expectedValue) {
		reflectType := reflect.TypeOf(expectedValue)
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}
		if typeName := reflectType.Name(); typeName != "" {
			t.Errorf("expected %s \"%v\", got \"%v\"", typeName, expectedValue, value)
			return
		}
		t.Errorf("expected \"%v\", got \"%v\"", expectedValue, value)
	}
}

func ErrorHasMessage(t *testing.T, err error, expectedMessage string) {
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != expectedMessage {
		t.Errorf("expected error \"%s\", got \"%s\"", expectedMessage, err.Error())
	}
}

func IsValidJSON(t *testing.T, value string, ptr any) {
	if err := json.Unmarshal([]byte(value), ptr); err != nil {
		t.Errorf(`failed to unmarshal "%s", got "%s"`, value, err.Error())
	}
}
