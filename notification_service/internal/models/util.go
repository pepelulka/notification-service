package models

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	Raw sql.NullString
}

// Custom JSON-serialization for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Raw.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Raw.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Raw.String = ""
		ns.Raw.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &ns.Raw.String)
	if err != nil {
		return err
	}
	ns.Raw.Valid = true
	return nil
}
