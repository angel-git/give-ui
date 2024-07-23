package models

type ViewItem struct {
	Id          string
	Name        string
	Type        string
	Description string
	Category    string
	MaxStock    int
}

type AllItems struct {
	Categories []string
	Items      []ViewItem
	Presets    []ViewItem
}
