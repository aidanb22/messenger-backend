package database

import (
	"context"
	"errors"
	"github.com/ablancas22/messenger-backend/models"
	"time"
)

// MessageService is used by the app to manage all Task related controllers and functionality
type MessageService struct {
	collection     DBCollection
	db             DBClient
	messageHandler *DBHandler[*messageModel]
	userHandler    *DBHandler[*userModel]
	groupHandler   *DBHandler[*groupModel]
}

// NewMessageService is an exported function used to initialize a new MessageService struct
func NewMessageService(db DBClient, tHandler *DBHandler[*messageModel], uHandler *DBHandler[*userModel], gHandler *DBHandler[*groupModel]) *MessageService {
	collection := db.GetCollection("messages")
	return &MessageService{collection, db, tHandler, uHandler, gHandler}
}

// checkLinkedRecords ensures the userId and groupId in the models.Task is correct
func (p *MessageService) checkMessageGroups(g *groupModel, u *userModel) error {
	gOutCh := make(chan *groupModel)
	gErrCh := make(chan error)
	uOutCh := make(chan *userModel)
	uErrCh := make(chan error)
	go func() {
		reG, err := p.groupHandler.FindOne(g)
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := p.userHandler.FindOne(u)
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
func (p *MessageService) checkMessageUsers(g *userModel, u *userModel) error {
	gOutCh := make(chan *userModel)
	gErrCh := make(chan error)
	uOutCh := make(chan *userModel)
	uErrCh := make(chan error)
	go func() {
		reG, err := p.userHandler.FindOne(g)
		gOutCh <- reG
		gErrCh <- err
	}()
	go func() {
		reU, err := p.userHandler.FindOne(u)
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

// TaskCreate is used to create a new user Task
func (p *MessageService) MessageCreate(g *models.Message) (*models.Message, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	gm, err := newMessageModel(g)
	if err != nil {
		return nil, err
	}
	if gm.Group {
		err = p.checkMessageGroups(&groupModel{Id: gm.ReceiverId}, &userModel{Id: gm.SenderId})
	} else {
		err = p.checkMessageUsers(&userModel{Id: gm.ReceiverId}, &userModel{Id: gm.SenderId})
	}
	if err != nil {
		return nil, err
	}
	gm, err = p.messageHandler.InsertOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// MessagesFind is used to find all Task docs in a MongoDB Collection
func (p *MessageService) MessagesFind(g *models.Message) ([]*models.Message, error) {
	var tasks []*models.Message
	tm, err := newMessageModel(g)
	if err != nil {
		return tasks, err
	}
	gms, err := p.messageHandler.FindMany(tm)
	if err != nil {
		return tasks, err
	}
	for _, gm := range gms {
		tasks = append(tasks, gm.toRoot())
	}
	return tasks, nil
}

// MessageFInd is used to find a specific Task doc
func (p *MessageService) MessageFind(g *models.Message) (*models.Message, error) {
	gm, err := newMessageModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.messageHandler.FindOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// Message delete is used to delete a Task doc
func (p *MessageService) MessageDelete(g *models.Message) (*models.Message, error) {
	gm, err := newMessageModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.messageHandler.DeleteOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// Messagedocinsert is used to insert a Task doc directly into mongodb for testing purposes
func (p *MessageService) MessageDocInsert(g *models.Message) (*models.Message, error) {
	insertTask, err := newMessageModel(g)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = p.collection.InsertOne(ctx, insertTask)
	if err != nil {
		return nil, err
	}
	return insertTask.toRoot(), nil
}
