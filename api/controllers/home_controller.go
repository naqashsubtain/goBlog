package controllers

import (
	"net/http"
	//"github.com/victorsteven/fullstack/api/responses"
	"github.com/naqash/goBlog/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome GO API")

}
