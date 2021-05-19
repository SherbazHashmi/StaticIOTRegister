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
	"net/http"
	"strconv"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	// reading in request
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user.Prepare()
	err = user.Validate("")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(s.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, userCreated.ID))
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}

	userRetrieved, err := user.FindUserByID(s.DB, uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, userRetrieved)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// pattern is to setup empty user entity, take the body data and convert it to a workable object
	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("[Error] Unauthorized"))
	}

	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}

	user.Prepare()
	err = user.Validate("update")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateUser(s.DB, uint32(uid))

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)

	}

	responses.JSON(w, http.StatusOK, updatedUser)
}


func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)

	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.DeleteUser(s.DB, uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))

	responses.JSON(w, http.StatusNoContent, "")
}