package models

type ViewItem struct {
	Id          string
	Name        string
	Type        string
	Description string
	Category    string
	MaxStock    int
}

type ViewPreset struct {
	Id           string
	Encyclopedia string
}

type AllItems struct {
	Categories    []string
	Items         map[string]ViewItem
	GlobalPresets []ViewPreset
}
