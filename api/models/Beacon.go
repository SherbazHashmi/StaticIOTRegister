package models

//TODO: Complete Beacon Object
// Current blocker: Unable to create beacon until organisation model is concretely defined

type Beacon struct {
	ID           uint64       `gorm:"primary;auto_increment" json:"id"`
	MacAddress   string       `gorm:"size:17; not null; unique" json:"mac_address"`
	Organization Organisation `json:"organisation"`
	IsActive     bool         `gorm:"default:false " json:"is_active"`
}