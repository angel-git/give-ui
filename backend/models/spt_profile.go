package models

type SPTProfileInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SPTProfile struct {
	Info SPTProfileInfo `json:"info"`
}
