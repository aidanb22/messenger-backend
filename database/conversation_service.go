package database

import (
	"context"
	"errors"
	"github.com/ablancas22/messenger-backend/models"
	"time"
)

type ConversationService struct {
	collection DBCollection
	db         DBClient
	handler    *DBHandler[*conversationModel]
}

func NewConversationService(db DBClient, handler *DBHandler[*conversationModel]) *ConversationService {
	collection := db.GetCollection("conversations")
	return &ConversationService{collection, db, handler}
}

// ConversationCreate is used to create a new user conversation
func (c *ConversationService) ConversationCreate(g *models.Conversation) (*models.Conversation, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	cm, err := newConversationModel(g)
	if err != nil {
		return nil, err
	}
	_, err = c.handler.FindOne(&conversationModel{Id: cm.Id, ParticipantsIds: cm.ParticipantsIds, Group: cm.Group})
	if err == nil {
		return nil, errors.New("conversation name exists")
	}
	cm, err = c.handler.InsertOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// ConversationsFind is used to find all group docs in a MongoDB Collection
func (c *ConversationService) ConversationsFind(g *models.Conversation) ([]*models.Conversation, error) {
	var groups []*models.Conversation
	//Todo: build filter using g
	cms, err := c.handler.FindMany(&conversationModel{})
	if err != nil {
		return groups, err
	}
	for _, cm := range cms {
		groups = append(groups, cm.toRoot())
	}
	return groups, nil
}

// ConversationFind is used to find a specific group doc
func (c *ConversationService) ConversationFind(g *models.Conversation) (*models.Conversation, error) {
	cm, err := newConversationModel(g)
	if err != nil {
		return nil, err
	}
	cm, err = c.handler.FindOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// ConversationDelete is used to delete a group doc
func (c *ConversationService) ConversationDelete(g *models.Conversation) (*models.Conversation, error) {
	cm, err := newConversationModel(g)
	if err != nil {
		return nil, err
	}
	cm, err = c.handler.DeleteOne(cm)
	if err != nil {
		return nil, err
	}
	return cm.toRoot(), err
}

// ConversationUpdate is used to update an existing group
func (c *ConversationService) ConversationUpdate(g *models.Conversation) (*models.Conversation, error) {
	var filter models.Conversation
	err := g.Validate("create")
	if err != nil {
		return nil, errors.New("missing valid query filter")
	}
	filter.Id = g.Id
	f, err := newConversationModel(&filter)
	if err != nil {
		return nil, err
	}
	cm, err := newConversationModel(g)
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

// ConversationDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (c *ConversationService) ConversationDocInsert(g *models.Conversation) (*models.Conversation, error) {
	insertGroup, err := newConversationModel(g)
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

/*

func (c *ConversationService) checkConversationGroups(g *groupModel, u *userModel) error {
	gOutCh := make(chan *groupModel)
	gErrCh := make(chan error)
	uOutCh := make(chan *userModel)
	uErrCh := make(chan error)
	go func() {
		reG, err := c.groupHandler.FindOne(g)
		gOutCh <- reG
		uErrCh <- err
	}()
	go func() {
		reU, err := c.userHandler.FindOne(u)
		uOutCh <- reU
		uErrCh <- err
	}()
	for i := 0; i < 4; i++ {
		select {
		case gOut := <-gOutCh:
			g = gOut
		case gErr := <-gErrCh:
			if gErr != nil {
				return errors.New("invalid group id")
			}
		case uOut := <-uOutCh:
			u = uOut
		case uErr := <-uErrCh:
			if uErr != nil {
				return errors.New("invalid user id")
			}
		}
	}
	return nil
}

func (c *ConversationService) checkConversationUsers(g *userModel, u *userModel) error {
	gOutCh := make(chan *userModel)
	gErrCh := make(chan error)
	uOutCh := make(chan *userModel)
	uErrCh := make(chan error)
	go func() {
		reG, err := c.userHandler.FindOne(g)
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := c.userHandler.FindOne(u)
		uOutCh <- reU
		uErrCh <- err
	}()
	for i := 0; i < 4; i++ {
		select {
		case gOut := <-gOutCh:
			g = gOut
		case gErr := <-gErrCh:
			if gErr != nil {
				return errors.New("invalid user id")
			}
		case uOut := <-uOutCh:
			u = uOut
		case uErr := <-uErrCh:
			if uErr != nil {
				return errors.New("invalid user id")
			}
		}
	}
	return nil
}

*/
