package models

type BsqQuestResponse struct {
	Data []BsgQuest `json:"data"`
}

type BsgQuest struct {
	Id        string `json:"_id"`
	QuestName string `json:"QuestName"`
	Image     string `json:"image"`
	Location  string `json:"location"`
	TraderId  string `json:"traderId"`
	Name      string `json:"name"`
}
