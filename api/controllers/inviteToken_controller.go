package controllers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naqash/goBlog/api/models"
	"github.com/naqash/goBlog/api/responses"
	"github.com/naqash/goBlog/api/utils/formaterror"
)

func (server *Server) CreateInviteToken(w http.ResponseWriter, r *http.Request) {

	inviteToken := models.InviteToken{}

	tokenCreated, err := inviteToken.SaveInvite(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusCreated, tokenCreated)
}

// func (server *Server) GetInviteToken(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	token := vars["token"]
// 	// if err != nil {
// 	// 	responses.ERROR(w, http.StatusBadRequest, err)
// 	// 	return
// 	// }
// 	InviteToken := models.InviteToken{}

// 	InviteTokenFound, err := InviteToken.FindInviteTokenByToken(server.DB, token)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return GetInviteToken
// 	}

// }

func (server *Server) SetMiddlewareInviteToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["inviteToken"]
		InviteToken := models.InviteToken{}
		InviteTokenFound, err := InviteToken.FindInviteTokenByToken(server.DB, token)
		if InviteTokenFound == nil || err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("invite token not Valid"))
			return
		}
		next(w, r)
	}
}
