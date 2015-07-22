package misc

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// Value is a value that drivers must be able to handle.
// It is either nil or an instance of one of these types:
//
//   int64
//   float64
//   bool
type Value interface{}

// Int64Wrapper is a wrapper type for dealing with inconsistent integer
// responses from OpenStack. In some cases, OpenStack will use "" to
// represent a null / non-existent integer, as opposed to 0 or null.
// This is analogous to the NullInt64 type in the database/sql package.
type Int64Wrapper struct {
	Int64 int64
	Valid bool
}

// Value implements the driver Valuer interface.
func (n *Int64Wrapper) Value() (Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// UnmarshalJSON converts the bytes to a Int64Wrapper
func (n *Int64Wrapper) UnmarshalJSON(data []byte) error {
	if data == nil {
		*n = Int64Wrapper{Valid: false}

		return nil
	}

	if data != nil && len(data) == 2 && data[0] == '"' && data[1] == '"' {
		*n = Int64Wrapper{Valid: false}

		return nil
	}

	var err error

	i, err := strconv.ParseInt(unwrappedValue(data), 10, 64)

	if err != nil {
		*n = Int64Wrapper{Valid: false}
		return err
	}

	*n = Int64Wrapper{
		Int64: i,
		Valid: true,
	}

	return nil
}

// MarshalJSON converts a Int64 wrapper to a []byte.
func (n Int64Wrapper) MarshalJSON() ([]byte, error) {
	var doc bytes.Buffer

	if !n.Valid {
		return []byte("\"\""), nil
	}

	err := json.NewEncoder(&doc).Encode(n.Int64)

	if err != nil {
		return nil, err
	}

	return doc.Bytes(), nil
}

// Float64Wrapper is a wrapper type for dealing with inconsistent float
// responses from OpenStack. In some cases, OpenStack will use "" to
// represent a null / non-existent float, as opposed to 0 or null.
// This is analogous to the NullFloat64 type in the database/sql package.
type Float64Wrapper struct {
	Float64 float64
	Valid   bool
}

// Value implements the driver Valuer interface.
func (n Float64Wrapper) Value() (Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

// UnmarshalJSON converts the bytes to a Float64Wrapper
func (n *Float64Wrapper) UnmarshalJSON(data []byte) error {
	if data == nil {
		*n = Float64Wrapper{Valid: false}

		return nil
	}

	if data != nil && len(data) == 2 && data[0] == '"' && data[1] == '"' {
		*n = Float64Wrapper{Valid: false}

		return nil
	}

	var err error

	f, err := strconv.ParseFloat(unwrappedValue(data), 64)

	if err != nil {
		*n = Float64Wrapper{Valid: false}
		return err
	}

	*n = Float64Wrapper{
		Float64: f,
		Valid:   true,
	}

	return nil
}

// MarshalJSON converts a Float64 wrapper to a []byte.
func (n Float64Wrapper) MarshalJSON() ([]byte, error) {
	var doc bytes.Buffer

	if !n.Valid {
		return []byte("\"\""), nil
	}

	err := json.NewEncoder(&doc).Encode(n.Float64)

	if err != nil {
		return nil, err
	}

	return doc.Bytes(), nil
}

// BoolWrapper is a wrapper type for dealing with inconsistent bool
// responses from OpenStack. In some cases, OpenStack will use "" to
// represent a null / non-existent boolean, as opposed to false or null.
// This is analogous to the NullBool type in the database/sql package.
type BoolWrapper struct {
	Bool  bool
	Valid bool
}

// Value implements the driver Valuer interface.
func (n BoolWrapper) Value() (Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// UnmarshalJSON converts the bytes to a BoolWrapper
func (n *BoolWrapper) UnmarshalJSON(data []byte) error {
	if data == nil {
		*n = BoolWrapper{Valid: false}

		return nil
	}

	if data != nil && len(data) == 2 && data[0] == '"' && data[1] == '"' {
		*n = BoolWrapper{Valid: false}

		return nil
	}

	var err error

	b, err := strconv.ParseBool(unwrappedValue(data))

	if err != nil {
		*n = BoolWrapper{Valid: false}
		return err
	}

	*n = BoolWrapper{
		Bool:  b,
		Valid: true,
	}

	return nil
}

// MarshalJSON converts a BoolWrapper to a []byte.
func (n BoolWrapper) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("\"\""), nil
	}

	if n.Bool {
		return []byte("\"True\""), nil
	}

	return []byte("\"False\""), nil
}

func unwrappedValue(data []byte) string {
	startIndex := 0
	endIndex := len(data)
	str := string(data)

	if strings.HasPrefix(str, "\"") {
		startIndex = startIndex + 1
	}

	if strings.HasSuffix(str, "\"") {
		endIndex = endIndex - 1
	}

	return str[startIndex:endIndex]
}
