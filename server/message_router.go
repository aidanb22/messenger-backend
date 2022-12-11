package server

import (
	"encoding/json"
	"github.com/ablancas22/messenger-backend/auth"
	"github.com/ablancas22/messenger-backend/models"
	"github.com/ablancas22/messenger-backend/services"
	"github.com/ablancas22/messenger-backend/utilities"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type messageRouter struct {
	aService *services.TokenService
	tService services.MessageService
}

// NewMessageRouter is a function that initializes a new groupRouter struct
func NewMessageRouter(router *mux.Router, a *services.TokenService, t services.MessageService) *mux.Router {
	gRouter := messageRouter{a, t}
	router.HandleFunc("/messages", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages", a.MemberTokenVerifyMiddleWare(gRouter.MessagesShow)).Methods("GET")
	router.HandleFunc("/messages", a.MemberTokenVerifyMiddleWare(gRouter.CreateMessage)).Methods("POST")
	router.HandleFunc("/messages/{messageId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages/{messageId}", a.MemberTokenVerifyMiddleWare(gRouter.MessageShow)).Methods("GET")
	router.HandleFunc("/messages/{messageId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteMessage)).Methods("DELETE")
	return router
}

// MessagesShow returns all messages to client
func (gr *messageRouter) MessagesShow(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	var filter models.Message
	filter.SenderID = tokenData.UserId
	filter.ReceiverID = tokenData.UserId
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	messages, err := gr.tService.MessagesFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.NewEncoder(w).Encode(messagesDTO{Messages: messages}); err != nil {
		return
	}
}

// CreateTask from a REST Request post body
func (gr *messageRouter) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &message); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	message.Id = utilities.GenerateObjectID()
	g, err := gr.tService.MessageCreate(&message)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(g); err != nil {
			return
		}
	}
}

// MessageShow shows a specific task
func (gr *messageRouter) MessageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	if messageId == "" || messageId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing messageId"})
		return
	}
	message, err := gr.tService.MessageFind(&models.Message{Id: messageId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		return
	}
	return
}

// DeleteMessage deletes a message
func (gr *messageRouter) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	if messageId == "" || messageId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	message, err := gr.tService.MessageDelete(&models.Message{Id: messageId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		return
	}
	return
}
