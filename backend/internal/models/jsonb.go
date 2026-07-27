package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB is a custom type to handle PostgreSQL JSONB / JSON columns
type JSONB json.RawMessage

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		str, okStr := value.(string)
		if okStr {
			*j = []byte(str)
			return nil
		}
		return errors.New("type assertion to []byte/string failed")
	}

	result := make([]byte, len(bytes))
	copy(result, bytes)
	*j = result
	return nil
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// MarshalJSON implements json.Marshaler
func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (j *JSONB) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], b...)
	return nil
}
