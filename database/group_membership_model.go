package database

import (
	"fmt"
	"github.com/ablancas22/messenger-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type groupMembershipModel struct {
	Id        primitive.ObjectID `bson:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id,omitempty"`
	GroupId   primitive.ObjectID `bson:"group_id,omitempty"`
	Admin     bool               `bson:"admin,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty"`
}

// newGroupModel initializes a new pointer to a groupModel struct from a pointer to a JSON Group struct
func newGroupMembershipModel(g *models.GroupMembership) (gm *groupMembershipModel, err error) {
	gm = &groupMembershipModel{
		Admin:     g.Admin,
		UpdatedAt: g.UpdatedAt,
		CreatedAt: g.CreatedAt,
		DeletedAt: g.DeletedAt,
	}
	if g.Id != "" && g.Id != "000000000000000000000000" {
		gm.Id, err = primitive.ObjectIDFromHex(g.Id)
	}
	if g.UserId != "" && g.UserId != "000000000000000000000000" {
		gm.UserId, err = primitive.ObjectIDFromHex(g.UserId)
	}
	if g.GroupId != "" && g.GroupId != "000000000000000000000000" {
		gm.GroupId, err = primitive.ObjectIDFromHex(g.GroupId)
	}
	return
}

// toRoot creates and return a new pointer to a Group JSON struct from a pointer to a BSON groupModel
func (g *groupMembershipModel) toRoot() *models.GroupMembership {
	return &models.GroupMembership{
		Id:        g.Id.Hex(),
		UserId:    g.UserId.Hex(),
		GroupId:   g.GroupId.Hex(),
		Admin:     g.Admin,
		UpdatedAt: g.UpdatedAt,
		CreatedAt: g.CreatedAt,
		DeletedAt: g.DeletedAt,
	}
}

func (g *groupMembershipModel) update(doc interface{}) (err error) {
	data, err := bsonMarshall(doc)
	if err != nil {
		return
	}
	gmm := groupMembershipModel{}
	err = bson.Unmarshal(data, &gmm)
	if len(gmm.Id.Hex()) > 0 && gmm.Id.Hex() != "000000000000000000000000" {
		g.Id = gmm.Id
	}

	return
}

// bsonLoad loads a bson doc into the userModel
func (g *groupMembershipModel) bsonLoad(doc bson.D) (err error) {
	bData, err := bsonMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, g)
	return err
}

// match compares an input bson doc and returns whether there's a match with the userModel
// TODO: Find a better way to write these model match methods
func (g *groupMembershipModel) match(doc interface{}) bool {
	data, err := bsonMarshall(doc)
	fmt.Println("\n\ncheckMatch", doc, data, err)
	if err != nil {
		return false
	}
	gmm := groupMembershipModel{}
	err = bson.Unmarshal(data, &gmm)
	if gmm.Id.Hex() != "" && gmm.Id.Hex() != "000000000000000000000000" {
		if g.Id == gmm.Id {
			return true
		}
		return false
	} else if gmm.GroupId.Hex() != "" && gmm.GroupId.Hex() != "000000000000000000000000" {
		if gmm.UserId.Hex() != "" && gmm.UserId.Hex() != "000000000000000000000000" {
			if g.UserId == gmm.UserId && g.GroupId == gmm.GroupId {
				return true
			}
		}
	}

	return false
}

// getID returns the unique identifier of the userModel
func (g *groupMembershipModel) getID() (id interface{}) {
	return g.Id
}

// addTimeStamps updates an userModel struct with a timestamp
func (g *groupMembershipModel) addTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	if newRecord {
		g.CreatedAt = currentTime
	}
}

// addObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (g *groupMembershipModel) addObjectID() {
	if g.Id.Hex() == "" || g.Id.Hex() == "000000000000000000000000" {
		g.Id = primitive.NewObjectID()
	}
}

// postProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (g *groupMembershipModel) postProcess() (err error) {
	//u.Password = ""
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// toDoc converts the bson userModel into a bson.D
func (g *groupMembershipModel) toDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(g)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// bsonFilter generates a bson filter for MongoDB queries from the userModel data
func (g *groupMembershipModel) bsonFilter() (doc bson.D, err error) {
	if g.Id.Hex() != "" && g.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", g.Id}}
	}
	if g.GroupId.Hex() != "" && g.GroupId.Hex() != "000000000000000000000000" {
		if g.UserId.Hex() != "" && g.UserId.Hex() != "000000000000000000000000" {
			doc = bson.D{{"user_id", g.UserId}, {"group_id", g.GroupId}}
		} else {
			doc = bson.D{{"group_id", g.GroupId}}
		}
	}
	fmt.Println("\n\nbsonFilter", doc)
	return
}

// bsonUpdate generates a bson update for MongoDB queries from the userModel data
func (g *groupMembershipModel) bsonUpdate() (doc bson.D, err error) {
	inner, err := g.toDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}
