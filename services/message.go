package services

import "github.com/ablancas22/messenger-backend/models"

// TaskService is an interface used to manage the relevant group doc controllers
type MessageService interface {
	MessageCreate(g *models.Message) (*models.Message, error)
	MessageFind(g *models.Message) (*models.Message, error)
	MessagesFind(g *models.Message) ([]*models.Message, error)
	MessageDelete(g *models.Message) (*models.Message, error)
	MessageDocInsert(g *models.Message) (*models.Message, error)
}
