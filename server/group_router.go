package server

import (
	"encoding/json"
	"fmt"
	"github.com/ablancas22/messenger-backend/auth"
	"github.com/ablancas22/messenger-backend/models"
	"github.com/ablancas22/messenger-backend/services"
	"github.com/ablancas22/messenger-backend/utilities"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type groupRouter struct {
	aService  *services.TokenService
	gService  services.GroupService
	uService  services.UserService
	gmService services.GroupMembershipService
}

// NewGroupRouter is a function that initializes a new groupRouter struct
func NewGroupRouter(router *mux.Router, a *services.TokenService, g services.GroupService, u services.UserService, gm services.GroupMembershipService) *mux.Router {
	gRouter := groupRouter{a, g, u, gm}
	router.HandleFunc("/groups", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups", a.AdminTokenVerifyMiddleWare(gRouter.GroupsShow)).Methods("GET")
	router.HandleFunc("/groups", a.AdminTokenVerifyMiddleWare(gRouter.CreateGroup)).Methods("POST")
	router.HandleFunc("/groups/{groupId}", utilities.HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups/{groupId}", a.AdminTokenVerifyMiddleWare(gRouter.GroupShow)).Methods("GET")
	router.HandleFunc("/groups", a.AdminTokenVerifyMiddleWare(gRouter.CreateGroup)).Methods("POST")
	router.HandleFunc("/groups/{groupId}", a.AdminTokenVerifyMiddleWare(gRouter.DeleteGroup)).Methods("DELETE")
	router.HandleFunc("/groups/{groupId}", a.AdminTokenVerifyMiddleWare(gRouter.ModifyGroup)).Methods("PATCH")
	router.HandleFunc("/groups/{groupId}/users", a.MemberTokenVerifyMiddleWare(gRouter.GetGroupUsers)).Methods("GET")
	router.HandleFunc("/groups/{groupId}/users/{userId}", a.MemberTokenVerifyMiddleWare(gRouter.DeleteGroupUser)).Methods("DELETE")
	router.HandleFunc("/groups/{groupId}/users/{userId}", a.MemberTokenVerifyMiddleWare(gRouter.AddGroupUser)).Methods("POST")
	return router
}

func (gr *groupRouter) AddGroupUser(w http.ResponseWriter, r *http.Request) {
	var groupMember *models.GroupMembership
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &groupMember); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	groupMember.Id = utilities.GenerateObjectID()
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	userId := vars["userId"]
	groupMember.GroupId = groupId
	groupMember.UserId = userId
	//1: use userId from tokenData to check for membership record
	gm, err := gr.gmService.GroupMembershipFind(&models.GroupMembership{UserId: tokenData.UserId, GroupId: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	if !gm.Admin {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: "unauthorized"})
		return
	}
	//2: if it is, then delete the groupMembershipDoc for that userId and groupId
	gm, err = gr.gmService.GroupMembershipCreate(groupMember)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(gm); err != nil {
		return
	}
}

func (gr *groupRouter) DeleteGroupUser(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	userId := vars["userId"]
	//1: use userId from tokenData to check for membership record
	gm, err := gr.gmService.GroupMembershipFind(&models.GroupMembership{UserId: tokenData.UserId, GroupId: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	if !gm.Admin {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: "unauthorized"})
		return
	}
	//2: if it is, then delete the groupMembershipDoc for that userId and groupId
	gm, err = gr.gmService.GroupMembershipDelete(&models.GroupMembership{UserId: userId, GroupId: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(gm); err != nil {
		return
	}
}

// GroupsShow returns all groups to client
func (gr *groupRouter) GroupsShow(w http.ResponseWriter, r *http.Request) {
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	groups, err := gr.gService.GroupsFind()
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.NewEncoder(w).Encode(groupsDTO{Groups: groups}); err != nil {
		return
	}
}

// CreateGroup from a REST Request post body
func (gr *groupRouter) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	authToken := r.Header.Get("Auth-Token")
	tokenData, err := auth.DecodeJWT(authToken)
	fmt.Println("\n\ntoken", tokenData)

	if err != nil {
		utilities.RespondWithError(w, http.StatusUnauthorized, utilities.JWTError{Message: err.Error()})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &group); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	group.Id = utilities.GenerateObjectID()
	fmt.Println("\n\npre_controllerA", group)
	g, err := gr.gService.GroupCreate(&group)
	fmt.Println("\n\npost_controllerB", g)

	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	fmt.Println("\n\npre_controllerB", g.Id, tokenData.UserId)
	gm, err := gr.gmService.GroupMembershipCreate(&models.GroupMembership{GroupId: g.Id, UserId: tokenData.UserId, Admin: true})
	fmt.Println("\n\npost_controllerB", gm, err)

	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(g); err != nil {
		return
	}

}

// ModifyGroup to update a group document
func (gr *groupRouter) ModifyGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	if groupId == "" || groupId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	var group models.Group
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = r.Body.Close(); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	if err = json.Unmarshal(body, &group); err != nil {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: err.Error()})
		return
	}
	group.Id = groupId
	g, err := gr.gService.GroupUpdate(&group)
	if err != nil {
		utilities.RespondWithError(w, http.StatusServiceUnavailable, utilities.JWTError{Message: err.Error()})
		return
	} else {
		w = utilities.SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err = json.NewEncoder(w).Encode(g); err != nil {
			return
		}
	}
}

// GroupShow shows a specific group
func (gr *groupRouter) GroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	if groupId == "" || groupId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	group, err := gr.gService.GroupFind(&models.Group{Id: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(group); err != nil {
		return
	}
	return
}

// DeleteGroup deletes a group
func (gr *groupRouter) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	if groupId == "" || groupId == "000000000000000000000000" {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	group, err := gr.gService.GroupDelete(&models.Group{Id: groupId})
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(group); err != nil {
		return
	}
	return
}

func (gr *groupRouter) GetGroupUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	groupId := vars["groupId"]
	if !utilities.CheckObjectID(groupId) {
		utilities.RespondWithError(w, http.StatusBadRequest, utilities.JWTError{Message: "missing groupId"})
		return
	}
	dto, err := gr.getGroupUsers(groupId)
	if err != nil {
		utilities.RespondWithError(w, http.StatusNotFound, utilities.JWTError{Message: err.Error()})
		return
	}
	dto.clean()
	w = utilities.SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		return
	}
	return
}
func (gr *groupRouter) getGroupUsers(groupId string) (*groupUsersDTO, error) {
	var dto groupUsersDTO
	gOutCh := make(chan *models.Group)
	gErrCh := make(chan error)
	uOutCh := make(chan []*models.GroupMembership)
	uErrCh := make(chan error)
	go func() {
		reG, err := gr.gService.GroupFind(&models.Group{Id: groupId})
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := gr.gmService.GroupMembershipsFind(&models.GroupMembership{GroupId: groupId})
		uOutCh <- reU
		uErrCh <- err
	}()
	for i := 0; i < 4; i++ {
		select {
		case gOut := <-gOutCh:
			dto.Group = gOut
		case gErr := <-gErrCh:
			if gErr != nil {
				return &dto, gErr
			}
		case uOut := <-uOutCh:
			dto.Users = uOut
		case uErr := <-uErrCh:
			if uErr != nil {
				return &dto, uErr
			}
		}
	}
	return &dto, nil
}
