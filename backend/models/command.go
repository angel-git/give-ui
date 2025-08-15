package models

type Command struct {
	Message string `json:"message"`
}

type FinishQuestCommand struct {
	QuestId string `json:"id"`
}
