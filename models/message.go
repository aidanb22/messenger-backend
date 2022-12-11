package models

import (
	"errors"
	"strings"
	"time"
)

// Message is a root struct that is used to store the json encoded data for/from a mongodb group doc.
type Message struct {
	Id             string    `json:"id,omitempty"`
	ConversationID string    `json:"conversation_id,omitempty"`
	SenderID       string    `json:"sender_id,omitempty"`
	ReceiverID     string    `json:"receiver_id,omitempty"`
	Content        string    `json:"content,omitempty"`
	ContentType    string    `json:"contentType,omitempty"`
	Group          bool      `json:"group,omitempty"`
	FileIds        string    `json:"file_ids"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
}

// checkID determines whether a specified ID is set or not
func (g *Message) checkID(chkId string) bool {
	switch chkId {
	case "id":
		if g.Id == "" || g.Id == "000000000000000000000000" {
			return false
		}
	case "sender_id":
		if g.SenderID == "" || g.SenderID == "000000000000000000000000" {
			return false
		}
	case "receiver_id":
		if g.ReceiverID == "" || g.ReceiverID == "000000000000000000000000" {
			return false
		}
	}
	return true
}

// Validate a Group for different scenarios such as loading TokenData, creating new Group, or updating a Group
func (g *Message) Validate(valCase string) (err error) {
	var missingFields []string
	switch valCase {
	case "create":
		if !g.checkID("sender_id") {
			missingFields = append(missingFields, "sender_id")
		}
		if !g.checkID("receiver_id") {
			missingFields = append(missingFields, "receiver_id")
		}
		if g.Content == "" {
			missingFields = append(missingFields, "content")
		}
	default:
		return errors.New("unrecognized validation case")
	}
	if len(missingFields) > 0 {
		return errors.New("missing the following message fields: " + strings.Join(missingFields, ", "))
	}
	return
}
