package database

// https://github.com/jmoiron/sqlx/blob/master/types/types.go

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

// JSONB is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONB additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONB json.RawMessage

var emptyJSON = JSONB("{}")

func IsEmptyJSON(val JSONB) bool {
	return reflect.DeepEqual(val, emptyJSON)
}

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSON, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONB: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONB) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return string(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONB) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSON
		} else {
			source = t
		}
	case nil:
		*j = emptyJSON
	default:
		return errors.New("Incompatible type for JSONB")
	}
	*j = JSONB(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONB) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = emptyJSON
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONB types.
func (j JSONB) String() string {
	return string(j)
}
