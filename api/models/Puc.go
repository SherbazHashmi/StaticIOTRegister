package models

import "time"

type Puc struct {
	ID                   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	ParentOrganisationID uint64
	ParentOrganisation   Organisation `json:"parent_organisation"`
	LastCheckedIn        time.Time    `gorm:"default: CURRENT_TIMESTAMP"`


}