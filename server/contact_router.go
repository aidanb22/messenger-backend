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

type contactRouter struct {
	aService *services.TokenService
	tService services.MessageService
	cService services.ContactService
	uService services.UserService
}

// NewContactRouter is a function that initializes a new groupRouter struct
func NewContactRouter(router *mux.Router, a *services.TokenService, t services.MessageService, c services.ContactService, u services.UserService) *mux.Router {
	gRouter := contactRouter{a, t, c, u}
	router.HandleFunc("/contacts", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/contacts", a.MemberTokenVerifyMiddleWare(gRouter.ContactsShow)).Methods("GET")
	router.HandleFunc("/contacts", a.MemberTokenVerifyMiddleWare(gRouter.CreateContact)).Methods("POST")
	router.HandleFunc("/contacts/{contactsId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/contacts/{contactsId}", a.MemberTokenVerifyMiddleWare(gRouter.ContactShow)).Methods("GET")
	router.HandleFunc("/contacts/{contactsId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteContact)).Methods("DELETE")
	return router
}

// ContactsShow returns all contacts to client
func (cr *contactRouter) ContactsShow(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	var filter models.Contact
	filter.RequesterId = tokenData.UserId
	filter.RecipientId = tokenData.UserId
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	contacts, err := cr.cService.ContactsFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.NewEncoder(w).Encode(contactsDTO{Contacts: contacts}); err != nil {
		return
	}
}

// CreateContact creates a contact from a REST Request post body
func (cr *contactRouter) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &contact); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	contact.Id = utilities.GenerateObjectID()
	g, err := cr.cService.ContactCreate(&contact)
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

// ContactShow shows a specific contact
func (cr *contactRouter) ContactShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	if messageId == "" || messageId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing messageId"})
		return
	}
	contact, err := cr.cService.ContactFind(&models.Contact{Id: messageId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(contact); err != nil {
		return
	}
	return
}

// DeletContact deletes a contact
func (cr *contactRouter) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactId := vars["contactId"]
	if contactId == "" || contactId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing conversationId"})
		return
	}
	conversation, err := cr.cService.ContactDelete(&models.Contact{Id: contactId})
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
