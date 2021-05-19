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

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post := models.Post{}

	err = json.Unmarshal(body, &post)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(" Unauthorized access"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(" Unauthorized access"))
	}

	postCreated, err := post.SavePost(s.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host,r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (s *Server) GetPosts(w http.ResponseWriter, _ *http.Request) {
	post := models.Post{}

	posts, err := post.FindAllPosts(s.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(s.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(" unauthorized access"))
		return
	}

	post := models.Post{}
	err = s.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized access"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized access"))
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	post := models.Post{}
	err = s.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("unauthorized"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	_, err = post.DeleteAPost(s.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")

}