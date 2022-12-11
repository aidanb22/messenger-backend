package database

import (
	"errors"
	"github.com/ablancas22/messenger-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// taskModel structures a group BSON document to save in a users collection
type messageModel struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	ConversationId primitive.ObjectID `bson:"conversation_id,omitempty"`
	SenderId       primitive.ObjectID `bson:"sender_id,omitempty"`
	ReceiverId     primitive.ObjectID `bson:"receiver_id,omitempty"`
	Content        string             `bson:"content,omitempty"`
	ContentType    string             `bson:"content_type,omitempty"`
	Group          bool               `bson:"group,omitempty"`
	FileIds        string             `bson:"file_ids"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty"`
	DeletedAt      time.Time          `bson:"deleted_at,omitempty"`
}

// newTaskModel initializes a new pointer to a userModel struct from a pointer to a JSON User struct
func newMessageModel(u *models.Message) (um *messageModel, err error) {
	um = &messageModel{

		Content:     u.Content,
		ContentType: u.ContentType,
		Group:       u.Group,
		UpdatedAt:   u.UpdatedAt,
		CreatedAt:   u.CreatedAt,
		DeletedAt:   u.DeletedAt,
	}
	if u.Id != "" && u.Id != "000000000000000000000000" {
		um.Id, err = primitive.ObjectIDFromHex(u.Id)
	}
	if u.SenderID != "" && u.SenderID != "000000000000000000000000" {
		um.SenderId, err = primitive.ObjectIDFromHex(u.SenderID)
	}
	if u.ReceiverID != "" && u.ReceiverID != "000000000000000000000000" {
		um.ReceiverId, err = primitive.ObjectIDFromHex(u.ReceiverID)
	}
	return
}

// update the userModel using an overwrite bson.D doc
func (u *messageModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	um := messageModel{}
	err = bson.Unmarshal(data, &um)
	if !um.UpdatedAt.IsZero() {
		u.UpdatedAt = um.UpdatedAt
	}
	return
}

// bsonLoad loads a bson doc into the userModel
func (u *messageModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, u)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (u *messageModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	um := messageModel{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if u.Id == um.Id {
			return true
		}
		return false
	}
	if um.SenderId.Hex() != "" && um.SenderId.Hex() != "000000000000000000000000" {
		if u.SenderId == um.SenderId {
			return true
		}
		return false
	}
	if um.ReceiverId.Hex() != "" && um.ReceiverId.Hex() != "000000000000000000000000" {
		if u.ReceiverId == um.ReceiverId {
			return true
		}
		return false
	}
	return false
}

// getID returns the unique identifier of the userModel
func (u *messageModel) getID() (id interface{}) {
	return u.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (u *messageModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	u.UpdatedAt = currentTime
	if newRecord {
		u.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (u *messageModel) addObjectID() {
	if u.Id.Hex() == "" || u.Id.Hex() == "000000000000000000000000" {
		u.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (u *messageModel) postProcess() (err error) {
	//u.Password = ""
	if u.SenderId.Hex() == "" {
		err = errors.New("user record does not have an email")
	}
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (u *messageModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(u)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (u *messageModel) bsonFilter() (doc bson.D, err error) {
	if u.Id.Hex() != "" && u.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", u.Id}}
	} else if u.ReceiverId.Hex() != "" && u.ReceiverId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"group_id", u.ReceiverId}}
	} else if u.SenderId.Hex() != "" && u.SenderId.Hex() != "000000000000000000000000" {
		doc = bson.D{{"user_id", u.SenderId}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (u *messageModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := u.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// toRoot creates and return a new pointer to a User JSON struct from a pointer to a BSON userModel
func (u *messageModel) toRoot() *models.Message {
	return &models.Message{
		Id:          u.Id.Hex(),
		SenderID:    u.SenderId.Hex(),
		ReceiverID:  u.ReceiverId.Hex(),
		Content:     u.Content,
		ContentType: u.ContentType,
		Group:       u.Group,
		UpdatedAt:   u.UpdatedAt,
		CreatedAt:   u.CreatedAt,
		DeletedAt:   u.DeletedAt,
	}
}
