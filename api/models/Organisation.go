package models

import "time"

type Organisation struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Region string `gorm:"size: 50; not null" json:"region"`
	Address string `json:"address"`
	AddressValidated bool `json:"address_validated"`
	AdministratorID uint64
	Administrator User `json:"administrator"`
	EntityName string `json:"entity_name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LastUsedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_used_at"`
}
