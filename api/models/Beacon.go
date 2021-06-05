package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Beacon struct {
	ID           uint64       `gorm:"primary;auto_increment" json:"id"`
	MacAddress   string       `gorm:"size:17; not null; unique" json:"mac_address"`
	OrganisationID uint64	`json:"organisation_id"`
	Organization Organisation `json:"organisation,omitempty"`
	IsRegistered     bool         `gorm:"default:false " json:"is_registered"`
	LastUpdated time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"last_updated"`
	RegisteredOn time.Time `json:"registered_on"`
}

type BeaconEventType struct {
	EventType string `gorm:"primary;not null; unique" json:"event_type"`
}


func (b *Beacon) BeforeSave() {
	b.LastUpdated = time.Now()
}

func (b *Beacon) Prepare() error {
	if b.MacAddress == "" {
		return errors.New("unable to update Beacon object as no Mac Address provided")
	}

	return nil
}

// Implementing CRUD for beacons

func (b *Beacon) FindBeaconByID(db *gorm.DB, id uint64) (*Beacon, error) {
	err := db.Debug().Model(&Beacon{}).Where("id = ?", id).Take(&b).Error
	if err != nil {
		return b, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(fmt.Sprintf(
			"beacon (ID: %d) not found", id))
	}
	return nil, err
}

func (b *Beacon) FindAllBeacons(db *gorm.DB) (*[]Beacon, error) {
	var beacons []Beacon
	err := db.Debug().Model(&Beacon{}).Limit(100).Find(&beacons).Error
	if err != nil {
		return &[]Beacon{}, err
	}
	return &beacons, err
}

func (b *Beacon) FindOrganisationBeacons(db *gorm.DB, organisationID uint64) (*[]Beacon, error) {
	var beacons []Beacon
	err := db.Debug().Model(&Beacon{}).Limit(100).Find("organisation_id = ?", organisationID).Take(&beacons).Error
	if err != nil {
		return &[]Beacon{}, err
	}
	return &beacons, nil
}

func (b *Beacon) UpdateBeacon(db *gorm.DB, uid uint64) error {
	// do pre-update operations
	b.BeforeSave()

	// update object
	db = db.Debug().Model(&Beacon{}).Where("id = ?", uid).Take(&Beacon{}).UpdateColumns(
		map[string]interface{}{
			"mac_address": b.MacAddress,
			"organisation": b.Organization,
			"organisation_id": b.OrganisationID,
			"is_registered": b.IsRegistered,
			"last_updated": time.Now(),
		})

	if db.Error != nil {
		return db.Error
	}

	// retrieve updated object for return
	err := db.Debug().Model(&Beacon{}).Where("id = ?", uid).Take(&b).Error

	if err != nil {
		return nil
	}
	return nil
}

func (b *Beacon) DeleteBeacon(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&Beacon{}).Where("id = ?", uid).Take(&Beacon{}).Delete(&Beacon{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (b *Beacon) RegisterBeacon(db *gorm.DB, organisation Organisation) error {
	err := b.changeBeaconRegistration(db, organisation, false)
	if err != nil {
		return err
	}
	return nil
}

func (b *Beacon) DeregisterBeacon(db *gorm.DB, organisation Organisation) error {
	err := b.changeBeaconRegistration(db, organisation, true)
	if err != nil {
		return err
	}
	return nil
}

// changeBeaconRegistration Associates a Beacon to an Organisation
func (b *Beacon) changeBeaconRegistration(db *gorm.DB, organisation Organisation, isRegistration bool) error {
	b.RegisteredOn = time.Now()
	b.LastUpdated = time.Now()

	if isRegistration {
		if organisation.ID == 0 {
			return errors.New("unable to change registration of unresolved organisation")
		}
		b.Organization = organisation
		b.OrganisationID = organisation.ID
	} else {
		b.IsRegistered = true
	}

	err := b.UpdateBeacon(db, b.ID)
	if err != nil {
		return errors.New("unable to register beacon")
	}
	return nil
}

// CheckInPUC Logs interaction of PUC with Beacon
func (b *Beacon) CheckInPUC(db *gorm.DB, pucID uint32) error{
	p := Puc{}
	// Check if PUC is registered to the beacon
	err := db.Debug().Model(&Puc{}).Where("id = ?", pucID).Take(&p).Error
	if err != nil {
		return errors.New(fmt.Sprintf("unable to find PUC with the following ID: %d", pucID))
	}
	p.LastBeaconCheckedInto = *b
	p.LastCheckedIn = time.Now()
	updatedPuc, err := p.UpdatePuc(db, pucID)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to update PUC with the following ID: %d", pucID))
	}

	if updatedPuc.LastCheckedIn != p.LastCheckedIn || updatedPuc.LastBeaconCheckedIntoID != p.LastBeaconCheckedIntoID{
		return errors.New(fmt.Sprintf("update operation on puc with the following ID: %d failed", pucID))
	}
	return nil
}