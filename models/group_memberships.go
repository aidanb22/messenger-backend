package models

import (
	"errors"
	"strings"
	"time"
)

type GroupMembership struct {
	Id        string    `json:"id,omitempty"`
	UserId    string    `json:"user_id,omitempty"`
	GroupId   string    `json:"group_id,omitempty"`
	Admin     bool      `json:"admin,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (g *GroupMembership) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	case "user_id":
		if g.UserId == "" || g.UserId == "000000000000000000000000" {
			return false
		}
	case "group_id":
		if g.GroupId == "" || g.GroupId == "000000000000000000000000" {
			return false
		}
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *GroupMembership) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if !g.checkID("group_id") {
			missingFields = append(missingFields, "group_id")
		}
		if !g.checkID("user_id") {
			missingFields = append(missingFields, "user_id")
		}
	case "update":
		if !g.checkID("id") {
			missingFields = append(missingFields, "id")
		}
	default:
		return errors.New("unrecognized validation case")
	}
	if len(missingFields) > 0 {
		return errors.New("missing the following group fields: " + strings.Join(missingFields, ", "))
	}
	return
}
