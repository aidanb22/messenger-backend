package database

import (
	"github.com/ablancas22/messenger-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type contactModel struct {
	Id          primitive.ObjectID `bson:"id,omitempty"`
	RequesterId string             `bson:"requester_id,omitempty"`
	RecipientId string             `bson:"recipient_id,omitempty"`
	Status      string             `bson:"status,omitempty"` //status can be pending, approved, rejected, blocked
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	DeletedAt   time.Time          `bson:"deleted_at,omitempty"`
}

// newGroupModel initializes a new pointer to a groupModel struct from a pointer to a JSON Group struct
func newContactModel(c *models.Contact) (cm *contactModel, err error) {
	cm = &contactModel{
		Status:    c.Status,
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
	if c.Id != "" && c.Id != "000000000000000000000000" {
		cm.Id, err = primitive.ObjectIDFromHex(c.Id)
	}
	return
}

// toRoot creates and return a new pointer to a Group JSON struct from a pointer to a BSON groupModel
func (c *contactModel) toRoot() *models.Contact {
	return &models.Contact{
		Id:        c.Id.Hex(),
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
}

func (c *contactModel) update(doc interface{}) (err error) {
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
func (c *contactModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, c)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (c *contactModel) match(doc interface{}) bool {
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
func (c *contactModel) getID() (id interface{}) {
	return c.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (c *contactModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	if newRecord {
		c.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (c *contactModel) addObjectID() {
	if c.Id.Hex() == "" || c.Id.Hex() == "000000000000000000000000" {
		c.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (c *contactModel) postProcess() (err error) {
	//u.Password = ""
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (c *contactModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(c)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (c *contactModel) bsonFilter() (doc bson.D, err error) {
	if c.Id.Hex() != "" && c.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", c.Id}}
	}
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (c *contactModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := c.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}
