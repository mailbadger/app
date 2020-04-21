package entities

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON entity represents a json raw message bytes.
type JSON json.RawMessage

// Value returns a string value.
func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

// Scan scans the value as []byte or string and appends the bytes to
// the raw message bytes.
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var s []byte
	switch v := value.(type) {
	case []byte:
		s = v
	case string:
		s = []byte(v)
	default:
		return errors.New("invalid Scan Source")
	}

	*j = append((*j)[0:0], s...)

	return nil
}

// MarshalJSON returns the raw message bytes, or "null"
// in case the json is nil.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON unmarshals the json data into the raw message bytes.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("nil value")
	}

	*j = append((*j)[0:0], data...)
	return nil
}

// IsNull checks if the bytes length equals zero, or
// the string has "null" value.
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// Equals compares two JSON values.
func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}
