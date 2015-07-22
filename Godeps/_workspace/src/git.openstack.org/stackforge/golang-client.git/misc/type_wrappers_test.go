package misc

import (
	"bytes"
	"encoding/json"
	"testing"
)

// Int64Wrapper

func TestMarshalInt64Wrapper(t *testing.T) {
	iw := &Int64Wrapper{
		Valid: true,
		Int64: 123456,
	}

	expected := "123456"

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestMarshalInt64WrapperInvalid(t *testing.T) {
	iw := &Int64Wrapper{
		Valid: false,
	}

	expected := "\"\""

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestUnmarshalInt64WrapperValid(t *testing.T) {
	var result Int64Wrapper

	decodeJSON(t, []byte("123456"), &result, false)

	expected := Int64Wrapper{
		Valid: true,
		Int64: 123456,
	}

	if result.Int64 != expected.Int64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalInt64WrapperInvalid(t *testing.T) {
	var result Int64Wrapper

	decodeJSON(t, []byte("\"\""), &result, false)

	expected := Int64Wrapper{
		Valid: false,
	}

	if result.Int64 != expected.Int64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalInt64WrapperInvalidOtherString(t *testing.T) {
	var result Int64Wrapper

	decodeJSON(t, []byte("\"string\""), &result, true)

	expected := Int64Wrapper{
		Valid: false,
	}

	if result.Int64 != expected.Int64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

// Float64Wrapper

func TestMarshalFloat64Wrapper(t *testing.T) {
	iw := &Float64Wrapper{
		Valid:   true,
		Float64: 123456.2,
	}

	expected := "123456.2"

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestMarshalFloat64WrapperInvalid(t *testing.T) {
	iw := &Float64Wrapper{
		Valid: false,
	}

	expected := "\"\""

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestUnmarshalFloat64WrapperValid(t *testing.T) {
	var result Float64Wrapper

	decodeJSON(t, []byte("123456.2"), &result, false)

	expected := Float64Wrapper{
		Valid:   true,
		Float64: 123456.2,
	}

	if result.Float64 != expected.Float64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalFloat64WrapperInvalid(t *testing.T) {
	var result Float64Wrapper

	decodeJSON(t, []byte("\"\""), &result, false)

	expected := Float64Wrapper{
		Valid: false,
	}

	if result.Float64 != expected.Float64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalFloat64WrapperInvalidOtherString(t *testing.T) {
	var result Float64Wrapper

	decodeJSON(t, []byte("\"string\""), &result, true)

	expected := Float64Wrapper{
		Valid: false,
	}

	if result.Float64 != expected.Float64 && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

// BoolWrapper

func TestMarshalBoolWrapper(t *testing.T) {
	iw := &BoolWrapper{
		Valid: true,
		Bool:  true,
	}

	expected := "\"True\""

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestMarshalBoolWrapperFalse(t *testing.T) {
	iw := &BoolWrapper{
		Valid: true,
		Bool:  false,
	}

	expected := "\"False\""

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestMarshalBoolWrapperInvalid(t *testing.T) {
	iw := &BoolWrapper{
		Valid: false,
	}

	expected := "\"\""

	resultString := string(encodeJSON(t, iw))

	if resultString != expected {
		t.Errorf("Expected: %#v\nfound %#v", expected, resultString)
	}
}

func TestUnmarshalBoolWrapperValid(t *testing.T) {
	var result BoolWrapper

	decodeJSON(t, []byte("\"True\""), &result, false)

	expected := BoolWrapper{
		Valid: true,
		Bool:  true,
	}

	if result.Bool != expected.Bool && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalBoolWrapperInvalid(t *testing.T) {
	var result BoolWrapper

	decodeJSON(t, []byte("\"\""), &result, false)

	expected := BoolWrapper{
		Valid: false,
	}

	if result.Bool != expected.Bool && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

func TestUnmarshalBoolWrapperInvalidOtherString(t *testing.T) {
	var result BoolWrapper

	decodeJSON(t, []byte("\"string\""), &result, true)

	expected := BoolWrapper{
		Valid: false,
	}

	if result.Bool != expected.Bool && result.Valid != expected.Valid {
		t.Errorf("Expected: %#v\nfound %#v", expected, result)
	}
}

// utility functions

func encodeJSON(t *testing.T, value interface{}) []byte {
	var src bytes.Buffer

	if err := json.NewEncoder(&src).Encode(value); err != nil {
		t.Error(err)
	}

	var dst bytes.Buffer
	if err := json.Compact(&dst, src.Bytes()); err != nil {
		t.Error(err)
	}

	return dst.Bytes()
}

func decodeJSON(t *testing.T, doc []byte, result interface{}, shouldFail bool) {
	if err := json.NewDecoder(bytes.NewReader(doc)).Decode(&result); err != nil && !shouldFail {
		t.Error(err)
	}
}
