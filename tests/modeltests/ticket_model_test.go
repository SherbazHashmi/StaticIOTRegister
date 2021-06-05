package modeltests

import (
	"github.com/SherbazHashmi/goblog/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllTickets(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatalf("Error refreshing user and ticket table %v\n", err)
	}
	_, _, err = seedUsersAndTickets()
	if err != nil {
		log.Fatalf("Error seeding user and tickettable %v\n", err)
	}
	tickets, err := ticketInstance.FindAllTickets(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the tickets: %v\n", err)
		return
	}
	assert.Equal(t, len(*tickets), 2)
}

func TestSaveTicket(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatalf("Error user and ticket refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newTicket := models.Ticket{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}
	savedTicket, err := newTicket.SaveTicket(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the ticket: %v\n", err)
		return
	}
	assert.Equal(t, newTicket.ID, savedTicket.ID)
	assert.Equal(t, newTicket.Title, savedTicket.Title)
	assert.Equal(t, newTicket.Content, savedTicket.Content)
	assert.Equal(t, newTicket.AuthorID, savedTicket.AuthorID)

}

func TestGetTicketByID(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatalf("Error refreshing user and ticket table: %v\n", err)
	}
	ticket, err := seedOneUserAndOneTicket()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundTicket, err := ticketInstance.FindTicketByID(server.DB, ticket.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundTicket.ID, ticket.ID)
	assert.Equal(t, foundTicket.Title, ticket.Title)
	assert.Equal(t, foundTicket.Content, ticket.Content)
}

func TestUpdateATicket(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatalf("Error refreshing user and ticket table: %v\n", err)
	}
	ticket, err := seedOneUserAndOneTicket()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	ticketUpdate := models.Ticket{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@gmail.com",
		AuthorID: ticket.AuthorID,
	}
	updatedTicket, err := ticketUpdate.UpdateATicket(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedTicket.ID, ticketUpdate.ID)
	assert.Equal(t, updatedTicket.Title, ticketUpdate.Title)
	assert.Equal(t, updatedTicket.Content, ticketUpdate.Content)
	assert.Equal(t, updatedTicket.AuthorID, ticketUpdate.AuthorID)
}

func TestDeleteATicket(t *testing.T) {

	err := refreshUserAndTicketTable()
	if err != nil {
		log.Fatalf("Error refreshing user and ticket table: %v\n", err)
	}
	ticket, err := seedOneUserAndOneTicket()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := ticketInstance.DeleteATicket(server.DB, ticket.ID, ticket.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
