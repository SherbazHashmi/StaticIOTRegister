package controllers

import (
	"fmt"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}


func (s *Server) Initialize(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName string) {
	var err error
	if dbDriver == "postgres" {
		dbURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		s.DB, err = gorm.Open(dbDriver, dbURL)

		if err != nil {
			fmt.Printf("Cannot connect to %s database", dbDriver)
			log.Fatal("[Error] ", err)
		} else {
			fmt.Printf("Database connection established.")
		}
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) // Database migration
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) Run(addr string){
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}