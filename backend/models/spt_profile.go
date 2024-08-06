package models

type SPTProfileInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Item struct {
	Id  string `json:"_id"`
	Tpl string `json:"_tpl"`
}

type WeaponBuild struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	Items []Item `json:"Items"`
}

type UserBuilds struct {
	WeaponBuilds []WeaponBuild `json:"weaponBuilds"`
}

type SPTProfile struct {
	Info       SPTProfileInfo `json:"info"`
	UserBuilds UserBuilds     `json:"userbuilds"`
}
