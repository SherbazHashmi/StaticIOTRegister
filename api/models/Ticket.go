package models

import (
	"errors"
	"github.com/SherbazHashmi/goblog/api/formaterror"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Ticket struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Ticket) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Ticket) Validate() error {
	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Ticket) SaveTicket(db *gorm.DB) (*Ticket, error) {
	var err error
	err = p.Validate()
	if err != nil {
		return &Ticket{}, errors.New(err.Error())
	}
	err = db.Debug().Model(&Ticket{}).Create(&p).Error

	if err != nil {
		return &Ticket{}, formaterror.FormatError(err.Error())
	}

	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return p, nil
}

func (p *Ticket) FindAllTickets(db *gorm.DB) (*[]Ticket, error) {
	var err error
	tickets := []Ticket{}

	err = db.Debug().Model(&Ticket{}).Limit(100).Find(&tickets).Error

	if err != nil {
		return &[]Ticket{}, err
	}

	if len(tickets) > 0 {
		for i, _ := range tickets {
			err = db.Debug().Model(&User{}).Where("id = ?", tickets[i].AuthorID).Take(&tickets[i].Author).Error
			if err != nil {
				return &[]Ticket{}, err
			}
		}
	}
	return &tickets, nil
}

func (p *Ticket) FindTicketByID(db *gorm.DB, pid uint64) (*Ticket, error) {
	var err error

	err = db.Debug().Model(&Ticket{}).Where("id = ?", pid).Take(&p).Error

	if err != nil {
		return &Ticket{}, err
	}

	return p, nil
}

func (p *Ticket) UpdateATicket(db *gorm.DB) (*Ticket, error) {
	var err error

	err = db.Debug().Model(&Ticket{}).Where("id = ?", p.ID).Updates(
		Ticket{
			Title: p.Title, Content: p.Content, UpdatedAt: time.Now(),
		}).Error

	if err != nil {
		return &Ticket{}, err
	}

	return p, nil
}

func (p *Ticket) DeleteATicket(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Ticket{}).Where("id = ? and author_id = ?", pid, uid).Take(&Ticket{}).Delete(&Ticket{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
