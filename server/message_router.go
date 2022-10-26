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

// NewTaskRouter is a function that initializes a new groupRouter struct
func NewTaskRouter(router *mux.Router, a *services.TokenService, t services.MessageService) *mux.Router {
	gRouter := messageRouter{a, t}
	router.HandleFunc("/messages", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages", a.MemberTokenVerifyMiddleWare(gRouter.MessagesShow)).Methods("GET")
	router.HandleFunc("/messages", a.MemberTokenVerifyMiddleWare(gRouter.CreateMessage)).Methods("POST")
	router.HandleFunc("/messages/{messageId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages/{messageId}", a.MemberTokenVerifyMiddleWare(gRouter.MessageShow)).Methods("GET")
	router.HandleFunc("/messages/{messageId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteMessage)).Methods("DELETE")
	return router
}

// TasksShow returns all tasks to client
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
	tasks, err := gr.tService.MessagesFind(&filter)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		return
	}
}

// CreateTask from a REST Request post body
func (gr *messageRouter) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var task models.Message
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &task); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	task.Id = utilities.GenerateObjectID()
	g, err := gr.tService.MessageCreate(&task)
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

// TaskShow shows a specific task
func (gr *messageRouter) MessageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["messageId"]
	if taskId == "" || taskId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	task, err := gr.tService.MessageFind(&models.Message{Id: taskId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		return
	}
	return
}

// DeleteTask deletes a task
func (gr *messageRouter) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["messageId"]
	if taskId == "" || taskId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing taskId"})
		return
	}
	task, err := gr.tService.MessageDelete(&models.Message{Id: taskId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		return
	}
	return
}
