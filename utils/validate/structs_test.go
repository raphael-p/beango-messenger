package validate

import (
	"fmt"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestStructFromJSON(t *testing.T) {
	type Nest2 struct {
		Leaf1 string
		Leaf2 string
	}

	type Nest1 struct {
		Regular Nest2
		Plain   JSONField[Nest2]
	}

	type ValidStruct struct {
		RegularWithJSON  string `json:"nameOnJSONTag"`
		RegularPlain     string
		RegularZeroable  bool `zeroable:"true"`
		Nested           Nest1
		Optional         JSONField[uint16]            `optional:"true"`
		OptionalNullable JSONField[map[string]string] `optional:"true" nullable:"true"`
		Nullable         JSONField[float32]           `nullable:"true"`
		Zeroable         JSONField[bool]              `zeroable:"true"`
		Plain            JSONField[[]int]
	}

	validStruct := ValidStruct{
		"hackney",
		"clerkenwell",
		true,
		Nest1{
			Nest2{"farringdon", "kew"},
			JSONField[Nest2]{Nest2{"putney", "camden"}, false, true},
		},
		JSONField[uint16]{256, false, true},
		JSONField[map[string]string]{
			map[string]string{"key1": "value1", "key2": "value2"}, false, true,
		},
		JSONField[float32]{3.14, false, true},
		JSONField[bool]{true, false, true},
		JSONField[[]int]{[]int{-100, 30, 255}, false, true},
	}

	t.Run("NormalWithPointer", func(t *testing.T) {
		missingFields, err := StructFromJSON(&validStruct)
		assert.IsNil(t, err)
		assert.HasLength(t, missingFields, 0)
	})

	t.Run("NormalWithStruct", func(t *testing.T) {
		missingFields, err := StructFromJSON(validStruct)
		assert.IsNil(t, err)
		assert.HasLength(t, missingFields, 0)
	})

	t.Run("NonStruct", func(t *testing.T) {
		testMap := map[string]string{"key1": "val1", "key2": "val2"}
		errString := "expected `value` to be a struct or its pointer, got %T"

		missingFields, err := StructFromJSON(testMap)
		assert.ErrorHasMessage(t, err, fmt.Sprintf(errString, testMap))
		assert.HasLength(t, missingFields, 0)

		missingFields, err = StructFromJSON(&testMap)
		assert.ErrorHasMessage(t, err, fmt.Sprintf(errString, &testMap))
		assert.HasLength(t, missingFields, 0)
	})

	t.Run("RegularFieldsZero", func(t *testing.T) {
		input := validStruct
		input.RegularWithJSON = ""
		input.RegularPlain = ""

		missingFields, err := StructFromJSON(&input)
		assert.IsNil(t, err)
		xMissingFields := []string{"nameOnJSONTag", "RegularPlain"}
		assert.DeepEquals(t, missingFields, xMissingFields)
	})

	t.Run("ZeroableRegularFieldZero", func(t *testing.T) {
		input := validStruct
		input.RegularZeroable = false

		// zeroable has no effect on regular fields
		missingFields, err := StructFromJSON(&input)
		assert.IsNil(t, err)
		assert.DeepEquals(t, missingFields, []string{"RegularZeroable"})
	})

	t.Run("NestedFieldLeavesZero", func(t *testing.T) {
		input := validStruct
		input.Nested.Regular.Leaf1 = ""
		input.Nested.Plain.Value.Leaf1 = ""

		missingFields, err := StructFromJSON(&input)
		assert.IsNil(t, err)
		xMissingFields := []string{"Nested.Regular.Leaf1", "Nested.Plain.Leaf1"}
		assert.DeepEquals(t, missingFields, xMissingFields)
	})

	t.Run("NestedRequiredFieldStructZero", func(t *testing.T) {
		input := validStruct
		input.Nested.Regular = Nest2{}
		input.Nested.Plain.Value = Nest2{}

		missingFields, err := StructFromJSON(&input)
		assert.IsNil(t, err)
		xMissingFields := []string{
			"Nested.Regular.Leaf1",
			"Nested.Regular.Leaf2",
			"Nested.Plain.Leaf1",
			"Nested.Plain.Leaf2",
		}
		assert.DeepEquals(t, missingFields, xMissingFields)
	})

	t.Run("JSONFieldTags", func(t *testing.T) {
		testCases := []struct {
			isSet, isNull, isZero bool
			xMissingFields        []string
		}{
			{false, false, false, []string{"Nullable", "Zeroable", "Plain"}},
			{false, false, true, []string{"Nullable", "Zeroable", "Plain"}},
			{false, true, false, []string{"Optional", "Nullable", "Zeroable", "Plain"}},
			{false, true, true, []string{"Optional", "Nullable", "Zeroable", "Plain"}},
			{true, true, false, []string{"Optional", "Zeroable", "Plain"}},
			{true, true, true, []string{"Optional", "Zeroable", "Plain"}},
			{true, false, true, []string{"Plain"}},
		}
		for _, testCase := range testCases {
			input := validStruct
			if !testCase.isSet {
				input.Optional.IsSet = false
				input.OptionalNullable.IsSet = false
				input.Nullable.IsSet = false
				input.Zeroable.IsSet = false
				input.Plain.IsSet = false
			}
			if testCase.isNull {
				input.Optional.IsNull = true
				input.OptionalNullable.IsNull = true
				input.Nullable.IsNull = true
				input.Zeroable.IsNull = true
				input.Plain.IsNull = true
			}
			if testCase.isZero {
				input.Optional.Value = 0
				input.OptionalNullable.Value = nil
				input.Nullable.Value = 0
				input.Zeroable.Value = false
				input.Plain.Value = nil
			}

			missingFields, err := StructFromJSON(&input)
			assert.IsNil(t, err)
			assert.DeepEquals(t, missingFields, testCase.xMissingFields)
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		jsonField := &JSONField[map[string]string]{}
		err := jsonField.UnmarshalJSON([]byte(`{"key": "value"}`))
		assert.IsNil(t, err)
		xValue := map[string]string{"key": "value"}
		assert.DeepEquals(t, *jsonField, JSONField[map[string]string]{xValue, false, true})
	})

	t.Run("NullValue", func(t *testing.T) {
		jsonField := &JSONField[string]{}
		err := jsonField.UnmarshalJSON([]byte("null"))
		assert.IsNil(t, err)
		assert.DeepEquals(t, *jsonField, JSONField[string]{"", true, true})
	})

	t.Run("NotSet", func(t *testing.T) {
		// if the field is absent from the JSON, UnmarshalJSON does not get called
		jsonField := &JSONField[float32]{}
		assert.DeepEquals(t, *jsonField, JSONField[float32]{0, false, false})
	})

	t.Run("WrongType", func(t *testing.T) {
		jsonField := &JSONField[int]{}
		err := jsonField.UnmarshalJSON([]byte(`definitely not an int`))
		assert.ErrorHasMessage(t, err, "invalid character 'd' looking for beginning of value")
		assert.DeepEquals(t, *jsonField, JSONField[int]{0, false, true})
	})
}

func TestPointerToStringStruct(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		validStruct := struct{ Field1, Field2 string }{}
		isValid := PointerToStringStruct(&validStruct)
		assert.Equals(t, isValid, true)
	})

	t.Run("Nil", func(t *testing.T) {
		isValid := PointerToStringStruct(nil)
		assert.Equals(t, isValid, false)
	})

	t.Run("NotPointer", func(t *testing.T) {
		invalidStruct := struct{ Field1, Field2 string }{}
		isValid := PointerToStringStruct(invalidStruct)
		assert.Equals(t, isValid, false)
	})

	t.Run("NonStringField", func(t *testing.T) {
		invalidStruct := &struct {
			Field1 int
			Field2 string
		}{}
		isValid := PointerToStringStruct(invalidStruct)
		assert.Equals(t, isValid, false)
	})

	t.Run("NotStruct", func(t *testing.T) {
		var nonStruct float32 = 3.14
		isValid := PointerToStringStruct(&nonStruct)
		assert.Equals(t, isValid, false)
	})
}
