package dbmodel

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal("")
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
		v.String = ""
	}
	return nil
}

func NewNullString(s string) NullString {
	if s == "" {
		return NullString{
			NullString: sql.NullString{String: "", Valid: false},
		}
	} else {
		return NullString{
			NullString: sql.NullString{String: s, Valid: true},
		}
	}
}
