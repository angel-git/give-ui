package models

type Bundle struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Items       map[string]uint16 `json:"items"`
}
