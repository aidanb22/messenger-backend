package database

import (
	"github.com/ablancas22/messenger-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type conversationModel struct {
	Id              primitive.ObjectID   `bson:"_id,omitempty"`
	ParticipantsIds []primitive.ObjectID `bson:"participants_ids,omitempty"`
	Group           bool                 `bson:"group,omitempty"` //if group, only ParticipantID is group id
	DeletedAt       time.Time            `bson:"deleted_at,omitempty"`
	UpdatedAt       time.Time            `bson:"updatedAt,omitempty"`
	CreatedAt       time.Time            `bson:"createdAt,omitempty"`
}

// newConversationModel initializes a new pointer to a userModel struct from a pointer to a JSON User struct
func newConversationModel(c *models.Conversation) (cm *conversationModel, err error) {
	cm = &conversationModel{
		Group:     c.Group,
		DeletedAt: c.DeletedAt,
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}
	if c.Id != "" && c.Id != "000000000000000000000000" {
		cm.Id, err = primitive.ObjectIDFromHex(c.Id)
	}
	var participantsIds []primitive.ObjectID
	for _, id := range c.ParticipantsIds {
		if id != "" && id != "000000000000000000000000" {
			oID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return cm, err
			}
			participantsIds = append(participantsIds, oID)
		}
	}
	cm.ParticipantsIds = participantsIds
	return
}

// toRoot creates and return a new pointer to a Group JSON struct from a pointer to a BSON groupModel
func (c *conversationModel) toRoot() *models.Conversation {
	return &models.Conversation{
		Id: c.Id.Hex(),
		//ParticipantsIds: c.ParticipantsIds.Hex(), //Todo: fix this
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
}

func (c *conversationModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	um := userModel{}
	err = bson.Unmarshal(data, &um)
	if len(um.Id.Hex()) > 0 && um.Id.Hex() != "000000000000000000000000" {
		c.Id = um.Id
	}

	return
}

// bsonLoad loads a bson doc into the userModel
func (c *conversationModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, c)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (c *conversationModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	if err != nil {
		return false
	}
	um := userModel{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if c.Id == um.Id {
			return true
		}
		return false
	}

	return false
}

// getID returns the unique identifier of the userModel
func (c *conversationModel) getID() (id interface{}) {
	return c.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (c *conversationModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	if newRecord {
		c.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (c *conversationModel) addObjectID() {
	if c.Id.Hex() == "" || c.Id.Hex() == "000000000000000000000000" {
		c.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (c *conversationModel) postProcess() (err error) {
	//u.Password = ""
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (c *conversationModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(c)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (c *conversationModel) bsonFilter() (doc bson.D, err error) {
	if c.Id.Hex() != "" && c.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", c.Id}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (c *conversationModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := c.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}
