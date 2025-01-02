package models

type Info struct {
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
	Info       Info       `json:"info"`
	UserBuilds UserBuilds `json:"userbuilds"`
	Characters Characters `json:"characters"`
}

type Characters struct {
	PMC PMC `json:"pmc"`
}

type PMC struct {
	InfoPMC     InfoPMC                  `json:"Info"`
	TradersInfo map[string]TraderProfile `json:"TradersInfo"`
	Skills      Skills                   `json:"Skills"`
}

type InfoPMC struct {
	Level int `json:"Level"`
}

type Skills struct {
	Common []SkillCommon `json:"Common"`
}

type SkillCommon struct {
	Id       string  `json:"Id"`
	Progress float32 `json:"Progress"`
}

type TraderProfile struct {
	SalesSum     float32 `json:"salesSum"`
	Standing     float32 `json:"standing"`
	LoyaltyLevel int     `json:"loyaltyLevel"`
}
