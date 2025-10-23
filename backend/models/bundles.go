package models

type Bundle struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Items       map[string]int `json:"items"`
}
