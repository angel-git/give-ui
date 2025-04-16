package models

import (
	"encoding/json"
	"fmt"
)

type Info struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ItemWithUpd struct {
	Id              string    `json:"_id"`
	Tpl             string    `json:"_tpl"`
	ParentID        *string   `json:"parentId"`
	SlotID          *string   `json:"slotId"`
	Location        *Location `json:"location"` // location can also be a number
	Upd             *Upd      `json:"upd"`
	ImageBase64     string
	SizeX           int
	SizeY           int
	BackgroundColor string
}

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
	Inventory   Inventory                `json:"Inventory"`
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

type Inventory struct {
	Items []ItemWithUpd `json:"items"`
	Stash string        `json:"stash"`
}

type Location struct {
	X int      `json:"x"`
	Y int      `json:"y"`
	R Rotation `json:"r"`
}

type Rotation string

// Implement UnmarshalJSON to support number or object
func (l *Location) UnmarshalJSON(data []byte) error {
	// Try unmarshalling into a number first (e.g. `5`)
	var dummyInt int
	if err := json.Unmarshal(data, &dummyInt); err == nil {
		// If it's just a number, ignore or assign defaults if needed
		*l = Location{} // or store dummyInt somewhere if meaningful
		return nil
	}

	// If not a number, try as full Location object
	type Alias Location
	var loc Alias
	if err := json.Unmarshal(data, &loc); err != nil {
		return err
	}
	*l = Location(loc)
	return nil
}

func (r *Rotation) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as number first
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		switch intVal {
		case 0:
			*r = "Horizontal"
		case 1:
			*r = "Vertical"
		default:
			*r = Rotation(fmt.Sprintf("Unknown(%d)", intVal))
		}
		return nil
	}

	// Try to unmarshal as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		*r = Rotation(strVal)
		return nil
	}

	return fmt.Errorf("invalid rotation format: %s", string(data))
}
