package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SherbazHashmi/goblog/api/auth"
	"github.com/SherbazHashmi/goblog/api/formaterror"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/SherbazHashmi/goblog/api/responses"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) CreateTicket(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ticket := models.Ticket{}

	err = json.Unmarshal(body, &ticket)
	log.Printf("author_id: %v", ticket.AuthorID)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized access"))
		return
	}
	//log.Printf("uid: %d, ticket: %v", uid, ticket)

	if ticket.AuthorID == 0 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Required Author"))
		return
	}
	if uid != ticket.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized access"))
		return
	}


	ticketCreated, err := ticket.SaveTicket(s.DB)
	if err != nil {

		if strings.Contains(err.Error(), "equired"){
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host,r.URL.Path, ticketCreated.ID))
	responses.JSON(w, http.StatusCreated, ticketCreated)
}

func (s *Server) GetTickets(w http.ResponseWriter, _ *http.Request) {
	ticket := models.Ticket{}

	tickets, err := ticket.FindAllTickets(s.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, tickets)
}

func (s *Server) GetTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	ticket := models.Ticket{}

	ticketReceived, err := ticket.FindTicketByID(s.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, ticketReceived)
}

func (s *Server) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the ticket id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the ticket exist
	ticket := models.Ticket{}
	err = s.DB.Debug().Model(models.Ticket{}).Where("id = ?", pid).Take(&ticket).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Ticket not found"))
		return
	}

	// If a user attempt to update a ticket not belonging to him
	if uid != ticket.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data ticketed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	ticketUpdate := models.Ticket{}
	err = json.Unmarshal(body, &ticketUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != ticketUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	ticketUpdate.Prepare()
	err = ticketUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ticketUpdate.ID = ticket.ID

	ticketUpdated, err := ticketUpdate.UpdateATicket(s.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, ticketUpdated)
}

func (s *Server) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	ticket := models.Ticket{}
	err = s.DB.Debug().Model(models.Ticket{}).Where("id = ?", pid).Take(&ticket).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("unauthorized"))
		return
	}

	if uid != ticket.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	_, err = ticket.DeleteATicket(s.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")

}