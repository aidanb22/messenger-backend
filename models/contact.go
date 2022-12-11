package models

import (
	"errors"
	"strings"
	"time"
)

type Contact struct {
	Id          string    `bson:"id,omitempty"`
	RequesterId string    `bson:"requester_id,omitempty"`
	RecipientId string    `bson:"recipient_id,omitempty"`
	Status      string    `bson:"status,omitempty"` //status can be pending, approved, rejected, blocked
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	DeletedAt   time.Time `bson:"deleted_at,omitempty"`
}

func (g *Contact) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	case "requester_id":
		if g.RequesterId == "" || g.RequesterId == "000000000000000000000000" {
			return false
		}
	case "recipient_id":
		if g.RecipientId == "" || g.RecipientId == "000000000000000000000000" {
			return false
		}
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *Contact) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if !g.checkID("requester_id") {
			missingFields = append(missingFields, "requester_id")
		}
		if !g.checkID("recipient_id") {
			missingFields = append(missingFields, "participants_ids")
		}
	case "update":
		if !g.checkID("id") {
			missingFields = append(missingFields, "id")
		}
		if g.Status == "" {
			missingFields = append(missingFields, "status")
		}

	default:
		return errors.New("unrecognized validation case")
	}
	if len(missingFields) > 0 {
		return errors.New("missing the following group fields: " + strings.Join(missingFields, ", "))
	}
	return
}
