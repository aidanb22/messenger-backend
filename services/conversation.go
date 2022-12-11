package services

import "github.com/ablancas22/messenger-backend/models"

type ConversationService interface {
	ConversationCreate(g *models.Conversation) (*models.Conversation, error)
	ConversationFind(g *models.Conversation) (*models.Conversation, error)
	ConversationsFind(g *models.Conversation) ([]*models.Conversation, error)
	ConversationDelete(g *models.Conversation) (*models.Conversation, error)
	ConversationDocInsert(g *models.Conversation) (*models.Conversation, error)
}
