package models

type UpdateSkillRequest struct {
	Skill    string `json:"skill"`
	Progress int    `json:"progress"`
}
