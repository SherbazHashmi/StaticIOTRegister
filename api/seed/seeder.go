package seed

import (
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/jinzhu/gorm"
	"log"
)

var users = []models.User {
	models.User{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var tickets = []models.Ticket {
	models.Ticket{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Ticket{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

var organisations = []models.Organisation {
	{
		Region: "Canberra",
		EntityName: "Ladomme Cafe",
	},
	{
		Region: "Canberra",
		EntityName: "Le Bon",
	},
	{
		Region: "Canberra",
		EntityName: "High Road",
	},
	{
		Region: "Canberra",
		EntityName: "Two Before Ten",
	},
	{
		Region: "Melbourne",
		EntityName: "Higher Ground",
	},

}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Ticket{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Ticket{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Ticket{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tickets[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Ticket{}).Create(&tickets[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tickets table: %v", err)
		}
	}
}
