package main

import (
	"encoding/json"
	"time"
)

type CDCRecord struct {
	Payload  Payload `json:"payload"`
	Schema   Schema  `json:"schema"`
	Name     string  `json:"name,omitempty"`
	Optional bool    `json:"optional,omitempty"`
	Type     string  `json:"type,omitempty"`
}

type Operation string

const (
	Read   Operation = "r"
	Create           = "c"
	Update           = "u"
	Delete           = "d"
)

type Payload struct {
	Before      string      `json:"before,omitempty"`
	After       string      `json:"after,omitempty"`
	Op          Operation   `json:"op,omitempty"`
	Patch       string      `json:"patch,omitempty"`
	Filter      string      `json:"filter,omitempty"`
	Source      interface{} `json:"source,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	TimestampMS time.Time   `json:"timestampMS"`
}

type SchemaField struct {
	Field    string `json:"field,omitempty"`
	Type     string `json:"type,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

type Schema struct {
	Type     string        `json:"type,omitempty"`
	Name     string        `json:"name,omitempty"`
	Optional bool          `json:"optional,omitempty"`
	Fields   []SchemaField `json:"fields,omitempty"`
}

// Unwrap takes a CDC Formatted JSON record and returns only the "after" payload (wrapped in the JSON with
// schema envelope).
func (r CDCRecord) Unwrap() (map[string]interface{}, error) {
	return nil, nil
}

func (r CDCRecord) Bytes() ([]byte, error) {
	return json.Marshal(r)
}

func parseCDCRecord(raw []byte) (CDCRecord, error) {
	var rec CDCRecord
	err := json.Unmarshal(raw, &rec)
	if err != nil {
		return CDCRecord{}, err
	}
	return rec, nil
}
