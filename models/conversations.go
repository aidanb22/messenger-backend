package models

import (
	"errors"
	"github.com/ablancas22/messenger-backend/utilities"
	"strings"
	"time"
)

type Conversation struct {
	Id              string    `json:"id,omitempty"`
	ParticipantsIds []string  `json:"participants_ids,omitempty"`
	Group           bool      `json:"group,omitempty"` //if group, only ParticipantID is group id
	DeletedAt       time.Time `json:"deleted_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

func (g *Conversation) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	case "participants_ids":
		for _, id := range g.ParticipantsIds {
			if id != "" && id != "000000000000000000000000" {
				return true
			}
		}
		return false
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *Conversation) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if !g.checkID("participants_ids") {
			missingFields = append(missingFields, "participants_ids")
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

func (g *Conversation) CheckParticipants(id string) bool {
	return utilities.IfStrInSlice(id, g.ParticipantsIds)
}
