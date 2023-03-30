package assert

import (
	"reflect"
	"testing"
)

func DeepEquals(t *testing.T, value any, expectedValue any) {
	if !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("expected %v, but got %v", expectedValue, value)
	}
}

func ErrorHasMessage(t *testing.T, err error, expectedMessage string) {
	if err.Error() != expectedMessage {
		t.Errorf("expected error %v, but got %v", err.Error(), expectedMessage)
	}
}
