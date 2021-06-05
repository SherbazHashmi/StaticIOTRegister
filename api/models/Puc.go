package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Puc struct {
	ID            uint64 `gorm:"primary_key;auto_increment" json:"id"`
	CurrentUserID uint64
	CurrentUser   *User `json:"current_user"`
	LastCheckedIn        time.Time    `gorm:"default: CURRENT_TIMESTAMP"`
	LastBeaconCheckedIntoID uint64
	LastBeaconCheckedInto Beacon `json:"last_beacon_checked_into"`
}

func (p *Puc) UpdatePuc(db *gorm.DB, uid uint32) (*Puc, error) {
	return nil,	nil
}