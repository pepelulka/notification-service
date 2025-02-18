package models

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	Raw sql.NullString
}

// Кастомная JSON-сериализация для NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Raw.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Raw.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	// Проверяем, пришел ли `null`
	if string(data) == "null" {
		ns.Raw.String = ""
		ns.Raw.Valid = false
		return nil
	}

	// Если не `null`, парсим строку
	err := json.Unmarshal(data, &ns.Raw.String)
	if err != nil {
		return err
	}
	ns.Raw.Valid = true
	return nil
}
