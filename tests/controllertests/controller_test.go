package controllertests

import (
	"github.com/SherbazHashmi/goblog/api/controllers"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/SherbazHashmi/goblog/tests"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var server = controllers.Server{}
var userInstance = models.User{}
var ticketInstance = models.Ticket{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))

	if err != nil {
		log.Fatalf("Unable to load environment varaibles.")
	}

	tests.SetupDatabase(&server)

	os.Exit(m.Run())
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error

	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error

	if err != nil {
		return err
	}

	log.Printf("successfully refreshed table")

	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func seedUsers() ([]models.User, error) {
	var err error
	users := []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndTicketTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Ticket{}).Error

	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Ticket{}).Error

	if err != nil {
		return err
	}

	log.Printf("successfully refreshed tabes")
	return nil
}

func seedOneUserAndOneTicket() (models.Ticket, error) {
	err := refreshUserAndTicketTable()

	if err != nil {
		return models.Ticket{}, err
	}

	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error

	if err != nil {
		return models.Ticket{}, err
	}

	ticket := models.Ticket{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}

	err = server.DB.Model(&models.Ticket{}).Create(&ticket).Error

	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func seedUsersAndTickets() ([]models.User, []models.Ticket, error) {
	var err error

	var users = []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var tickets = []models.Ticket{
		models.Ticket{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		models.Ticket{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tickets[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Ticket{}).Create(&tickets[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tickets table: %v", err)
		}
	}
	return users, tickets, nil
}
