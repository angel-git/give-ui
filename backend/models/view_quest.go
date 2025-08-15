package models

type ViewQuest struct {
	QID        string
	Name       string
	Location   string
	Trader     string
	Image      string
	Conditions []ViewCondition
}

type ViewCondition struct {
	Name        string
	IsCompleted bool
}
