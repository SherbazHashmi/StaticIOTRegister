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
var postInstance = models.Post{}

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

func refreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error

	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error

	if err != nil {
		return err
	}

	log.Printf("successfully refreshed tabes")
	return nil
}

func seedOneUserAndOnePost() (models.Post, error) {
	err := refreshUserAndPostTable()

	if err != nil {
		return models.Post{}, err
	}

	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error

	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}

	err = server.DB.Model(&models.Post{}).Create(&post).Error

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Post, error) {
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
	var posts = []models.Post{
		models.Post{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		models.Post{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
