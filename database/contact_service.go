package database

import (
	"context"
	"errors"
	"github.com/ablancas22/messenger-backend/models"
	"time"
)

// GroupService is used by the app to manage all group related controllers and functionality
type ContactService struct {
	collection DBCollection
	db         DBClient
	handler    *DBHandler[*contactModel]
}

// NewGroupService is an exported function used to initialize a new GroupService struct
func NewContactService(db DBClient, handler *DBHandler[*contactModel]) *ContactService {
	collection := db.GetCollection("group_memberships")
	return &ContactService{collection, db, handler}
}

// GroupCreate is used to create a new user group
func (c *ContactService) ContactCreate(g *models.Contact) (*models.Contact, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	cm, err := newContactModel(g)
	if err != nil {
		return nil, err
	}
	_, err = c.handler.FindOne(&contactModel{Id: cm.Id, RequesterId: cm.RequesterId, RecipientId: cm.RecipientId, Status: cm.Status})
	if err == nil {
		return nil, errors.New("group name exists")
	}
	cm, err = c.handler.InsertOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// GroupsFind is used to find all group docs in a MongoDB Collection
func (c *ContactService) ContactsFind(g *models.Contact) ([]*models.Contact, error) {
	var groups []*models.Contact
	//Todo: build filter using g
	cms, err := c.handler.FindMany(&contactModel{})
	if err != nil {
		return groups, err
	}
	for _, cm := range cms {
		groups = append(groups, cm.toRoot())
	}
	return groups, nil
}

// GroupFind is used to find a specific group doc
func (c *ContactService) ContactFind(g *models.Contact) (*models.Contact, error) {
	cm, err := newContactModel(g)
	if err != nil {
		return nil, err
	}
	cm, err = c.handler.FindOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// GroupDelete is used to delete a group doc
func (c *ContactService) ContactDelete(g *models.Contact) (*models.Contact, error) {
	cm, err := newContactModel(g)
	if err != nil {
		return nil, err
	}
	cm, err = c.handler.DeleteOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// GroupUpdate is used to update an existing group
func (c *ContactService) ContactUpdate(g *models.Contact) (*models.Contact, error) {
	var filter models.Contact
	err := g.Validate("create")
	if err != nil {
		return nil, errors.New("missing valid query filter")
	}
	filter.Id = g.Id
	f, err := newContactModel(&filter)
	if err != nil {
		return nil, err
	}
	cm, err := newContactModel(g)
	if err != nil {
		return nil, err
	}
	_, groupErr := c.handler.FindOne(f)
	if groupErr != nil {
		return nil, errors.New("group not found")
	}
	cm, err = c.handler.UpdateOne(f, cm)
	return cm.toRoot(), err
}

// GroupDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (c *ContactService) ContactDocInsert(g *models.Contact) (*models.Contact, error) {
	insertGroup, err := newContactModel(g)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = c.collection.InsertOne(ctx, insertGroup)
	if err != nil {
		return nil, err
	}
	return insertGroup.toRoot(), nil
}
