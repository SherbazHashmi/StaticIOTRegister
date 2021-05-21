package modeltests

import (
	"fmt"
	"github.com/SherbazHashmi/goblog/api/controllers"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

// convention
func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("error getting env %v\n", err)
	}

	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	testDbDriver := os.Getenv("TEST_DB_DRIVER")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PORT"),
	)

	server.DB, err = gorm.Open(testDbDriver, dsn)

	if err != nil {
		fmt.Printf("Cannot connect to the database with the following properties\n%s\n", dsn)
		log.Fatal("[ERROR]", err)
		return
	}

	fmt.Printf("test database connection established")
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

	log.Printf("successfully refreshed tables")
	return nil
}

func seedOneUser() (models.User, error) {
	refreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error

	if err != nil {
		log.Fatalf("cannot seed users table %v", err)
		return user, err
	}

	return user, nil
}

func seedUsers() []error {
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

	var errors []error

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Printf("unable to create user %s", users[i].Email)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func refreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error

	if err != nil {
		log.Fatalf("[Error] Unable to drop tables for testing, %v", err)
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error

	if err != nil {
		log.Fatalf("[Error] Unable to migrate tables for testing, %v", err)
		return err
	}

	log.Printf("successfully refreshed tables for testing")
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
		log.Fatalf("[ERR] Unable to create user for testing %v", err)
		return models.Post{}, err
	}

	post := models.Post{
		Title:    "Best title ever",
		Content:  "Awesome content, keep coming back",
		AuthorID: user.ID,
	}

	err = server.DB.Model(&models.Post{}).Create(&post).Error

	if err != nil {
		log.Fatalf("[ERR] Unable to create post\n %v", err)
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
