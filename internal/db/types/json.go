package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type NullJSON struct {
	RawMessage json.RawMessage
	Valid      bool
}

func (j *NullJSON) Scan(src any) error {
	if src == nil {
		j.RawMessage = nil
		j.Valid = false
		return nil
	}

	switch v := src.(type) {
	case string:
		j.RawMessage = append(j.RawMessage[:0], v...)
	case []byte:
		j.RawMessage = append(j.RawMessage[:0], v...)
	default:
		return fmt.Errorf("cannot scan %T into NullJSON", src)
	}

	j.Valid = true
	return nil
}

func (j NullJSON) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return string(j.RawMessage), nil
}
