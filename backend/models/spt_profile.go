package models

type Info struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ItemWithUpd struct {
	Id     string `json:"_id"`
	Tpl    string `json:"_tpl"`
	ParentID    *string `json:"parentId"`
	SlotID      *string `json:"slotId"`
	Upd         *Upd    `json:"upd"`
	ImageBase64 string
}

//templ GearPreset(equipmentBuild models.EquipmentBuild) {
//{{ armorItem := findBySlotId(equipmentBuild.Items, "ArmorVest") }}
//{{ tacticalVest := findBySlotId(equipmentBuild.Items, "TacticalVest") }}
//{{ Earpiece := findBySlotId(equipmentBuild.Items, "Earpiece") }}
//{{ Eyewear := findBySlotId(equipmentBuild.Items, "Eyewear") }}
//{{ Holster := findBySlotId(equipmentBuild.Items, "Holster") }}
//{{ FaceCover := findBySlotId(equipmentBuild.Items, "FaceCover") }}
//{{ firstPrimaryWeapon := findBySlotId(equipmentBuild.Items, "FirstPrimaryWeapon") }}
//{{ helmet := findBySlotId(equipmentBuild.Items, "Headwear") }}
//<div>
//<h1>Armor!!!</h1>
//<img alt="item" st-yle="max-height: 200px" src={ fmt.Sprintf("https://assets.tarkov.dev/%s-base-image.png", armorItem.Tpl) }/>
//<h1>Primary weapon</h1>
//<img alt="item" style="max-height: 200px" src={ fmt.Sprintf("https://assets.tarkov.dev/%s-base-image.png", firstPrimaryWeapon.Tpl) }/>
//<h1>Helmet</h1>
//<img alt="item" style="max-height: 200px" src={ fmt.Sprintf("https://assets.tarkov.dev/%s-base-image.png", helmet.Tpl) }/>
//</div>
//<div>
//for _, item := range equipmentBuild.Items {
//<div>{ item.Id } - { item.SlotId }</div>
//<img alt="item" style="max-height: 200px" src={ fmt.Sprintf("https://assets.tarkov.dev/%s-base-image.png", item.Tpl) }/>
//}
//</div>
//}

type WeaponBuild struct {
	Id          string         `json:"Id"`
	Name        string         `json:"Name"`
	Root        string         `json:"Root"`
	Items       *[]ItemWithUpd `json:"Items"`
	ImageBase64 string
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

type EquipmentBuild struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	Items []Item `json:"Items"`
}

type UserBuilds struct {
	WeaponBuilds    []WeaponBuild    `json:"weaponBuilds"`
	MagazineBuilds  []MagazineBuild  `json:"magazineBuilds"`
	EquipmentBuilds []EquipmentBuild `json:"equipmentBuilds"`
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
	Level    int    `json:"Level"`
	Nickname string `json:"Nickname"`
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

type Upd struct {
	Togglable         *Togglable `json:"Togglable"`
	Foldable          *Foldable  `json:"Foldable"`
	StackObjectsCount int        `json:"StackObjectsCount"`
}

type Togglable struct {
	On bool `json:"on"`
}

type Foldable struct {
	Folded bool `json:"folded"`
}
