package server

import (
	"errors"
	"fmt"
	"github.com/ablancas22/messenger-backend/models"
)

/*
================ User DTOs ==================
*/

// updatePassword is used when updating a user password
type updatePassword struct {
	NewPassword     string `json:"new_password"`
	CurrentPassword string `json:"current_password"`
}

// userSignIn is used when updating a user password
type userSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// toUser converts userSignIn DTO to a user
func (u *userSignIn) toUser() (*models.User, error) {
	if u.Email == "" {
		return &models.User{}, errors.New("missing user email")
	}
	if u.Password == "" {
		return &models.User{}, errors.New("missing user password")
	}
	return &models.User{
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

// usersDTO is used when returning a slice of User
type usersDTO struct {
	Users []*models.User `json:"users"`
}

// clean ensures the users in the usersDTO have no passwords set
func (u *usersDTO) clean() {
	for i, _ := range u.Users {
		u.Users[i].Password = ""
	}
}

// userMessagesDTO is used when returning user with associated tasks
type userMessagesDTO struct {
	User     *models.User      `json:"user"`
	Messages []*models.Message `json:"messages"`
}

// clean ensures the users in the userTasksDTO have password set
func (u *userMessagesDTO) clean() {
	u.User.Password = ""
}

/*
================ Group DTOs ==================
*/

// groupsDTO is used when returning a slice of Group
type groupsDTO struct {
	Groups []*models.Group `json:"groups"`
}

// groupUsersDTO is used when returning a group with its associated users
type groupUsersDTO struct {
	Group *models.Group             `json:"group"`
	Users []*models.GroupMembership `json:"users"`
}

// clean ensures the users in the groupUsersDTO have no passwords set
func (u *groupUsersDTO) clean() {
	for i, _ := range u.Users {
		//u.Users[i].Password = ""
		fmt.Println(i)
	}
}

// groupMessagesDTO is used when returning a group with its associated tasks
type groupMessagesDTO struct {
	Group    *models.Group     `json:"group"`
	Messages []*models.Message `json:"messages"`
}

/*
================ Messages DTOs ==================
*/

// messagesDTO is used when returning a slice of Task
type messagesDTO struct {
	Messages []*models.Message `json:"messages"`
}

/*
================ Conversations DTOs ==================
*/

// messagesDTO is used when returning a slice of Task
type conversationsDTO struct {
	Conversations []*models.Conversation `json:"conversations"`
}

/*
================ Contacts DTOs ==================
*/

// messagesDTO is used when returning a slice of Task
type contactsDTO struct {
	Contacts []*models.Contact `json:"contacts"`
}

/*
================ GroupMemberships DTOs ==================
*/

// messagesDTO is used when returning a slice of Task
type groupMembershipsDTO struct {
	GroupMemberships []*models.GroupMembership `json:"groupMemberships"`
}
