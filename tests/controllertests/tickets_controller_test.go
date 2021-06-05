package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateTicket(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title already taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized access",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized access",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized access",
		},
	}
	for _, v := range samples {
		log.Printf("---\nMaking Request\n%v", v)
		req, err := http.NewRequest("POST", "/tickets", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("unable to start http ticket request error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateTicket)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Printf("\n%v", rr.Body.String())
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetTickets(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndTickets()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/tickets", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetTickets)
	handler.ServeHTTP(rr, req)

	var tickets []models.Ticket
	err = json.Unmarshal([]byte(rr.Body.String()), &tickets)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(tickets), 2)
}
func TestGetTicketByID(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatal(err)
	}
	ticket, err := seedOneUserAndOneTicket()
	if err != nil {
		log.Fatal(err)
	}
	ticketSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(ticket.ID)),
			statusCode: 200,
			title:      ticket.Title,
			content:    ticket.Content,
			author_id:  ticket.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range ticketSample {

		req, err := http.NewRequest("GET", "/tickets", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetTicket)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, ticket.Title, responseMap["title"])
			assert.Equal(t, ticket.Content, responseMap["content"])
			assert.Equal(t, float64(ticket.AuthorID), responseMap["author_id"]) //the response author id is float64
		}
	}
}

func TestUpdateTicket(t *testing.T) {

	var TicketUserEmail, TicketUserPassword string
	var AuthTicketAuthorID uint32
	var AuthTicketID uint64

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tickets, err := seedUsersAndTickets()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		TicketUserEmail = user.Email
		TicketUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(TicketUserEmail, TicketUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first ticket
	for _, ticket := range tickets {
		if ticket.ID == 2 {
			continue
		}
		AuthTicketID = ticket.ID
		AuthTicketAuthorID = ticket.AuthorID
	}
	// fmt.Printf("this is the auth ticket: %v\n", AuthTicketID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"The updated ticket", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			title:        "The updated ticket",
			content:      "This is the updated content",
			author_id:    AuthTicketAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to ticket 2, and title must be unique
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title already taken",
		},
		{
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthTicketID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/tickets", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateTicket)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			log.Printf("valid response: \n%v", v)

			log.Printf("actual response: \n%v", responseMap)
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteTicket(t *testing.T) {

	var TicketUserEmail, TicketUserPassword string
	var TicketUserID uint32
	var AuthTicketID uint64

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tickets, err := seedUsersAndTickets()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		TicketUserEmail = user.Email
		TicketUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(TicketUserEmail, TicketUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second ticket
	for _, ticket := range tickets {
		if ticket.ID == 1 {
			continue
		}
		AuthTicketID = ticket.ID
		TicketUserID = ticket.AuthorID
	}
	ticketSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTicketID)),
			author_id:    TicketUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthTicketID)),
			author_id:    TicketUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthTicketID)),
			author_id:    TicketUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "unauthorized",
		},
	}
	for _, v := range ticketSample {

		req, _ := http.NewRequest("GET", "/tickets", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteTicket)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}