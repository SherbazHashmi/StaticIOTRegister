package controllers

import (
	"github.com/SherbazHashmi/goblog/api/responses"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "v0")
}