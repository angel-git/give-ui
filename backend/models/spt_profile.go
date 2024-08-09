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

type MagazineBuild struct {
	Id      string              `json:"Id"`
	Name    string              `json:"Name"`
	Caliber string              `json:"Caliber"`
	Items   []MagazineBuildItem `json:"Items"`
}

type MagazineBuildItem struct {
	TemplateId string `json:"TemplateId"`
	Count      int    `json:"Count"`
}

type UserBuilds struct {
	WeaponBuilds   []WeaponBuild   `json:"weaponBuilds"`
	MagazineBuilds []MagazineBuild `json:"magazineBuilds"`
}

type SPTProfile struct {
	Info       SPTProfileInfo `json:"info"`
	UserBuilds UserBuilds     `json:"userbuilds"`
}
