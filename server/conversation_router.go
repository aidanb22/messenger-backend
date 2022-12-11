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

type conversationRouter struct {
	aService *services.TokenService
	tService services.MessageService
	cService services.ConversationService
	uService services.UserService
}

//NewConversationRouter is a function that initializes a new groupRouter struct
func NewConversationRouter(router *mux.Router, a *services.TokenService, t services.MessageService, c services.ConversationService, u services.UserService) *mux.Router {
	gRouter := conversationRouter{a, t, c, u}
	router.HandleFunc("/conversations", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/conversations", a.MemberTokenVerifyMiddleWare(gRouter.ConversationsShow)).Methods("GET")
	router.HandleFunc("/conversations", a.MemberTokenVerifyMiddleWare(gRouter.CreateConversation)).Methods("POST")
	router.HandleFunc("/conversations/{conversationId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/conversations/{conversationsId}", a.MemberTokenVerifyMiddleWare(gRouter.ConversationShow)).Methods("GET")
	router.HandleFunc("/conversations/{conversationsId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteConversation)).Methods("DELETE")
	return router
}

// ConversationsShow returns all conversations to client
func (gr *conversationRouter) ConversationsShow(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	var filter models.Conversation
	filter.ParticipantsIds = append(filter.ParticipantsIds, tokenData.UserId)

	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	conversations, err := gr.cService.ConversationsFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.NewEncoder(w).Encode(conversationsDTO{Conversations: conversations}); err != nil {
		return
	}
}

// CreateConversation from a REST Request post body
func (gr *conversationRouter) CreateConversation(w http.ResponseWriter, r *http.Request) {
	var conversation models.Conversation
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &conversation); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	conversation.Id = utilities.GenerateObjectID()
	g, err := gr.cService.ConversationCreate(&conversation)
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

// ConversationShow shows a specific conversation
func (gr *conversationRouter) ConversationShow(w http.ResponseWriter, r *http.Request) {
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
func (cr *conversationRouter) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationId := vars["conversationId"]
	if conversationId == "" || conversationId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing conversationId"})
		return
	}
	conversation, err := cr.cService.ConversationDelete(&models.Conversation{Id: conversationId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(conversation); err != nil {
		return
	}
	return
}
