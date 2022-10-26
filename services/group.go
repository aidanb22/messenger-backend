package services

import "github.com/ablancas22/messenger-backend/models"

// GroupService is an interface used to manage the relevant group doc controllers
type GroupService interface {
	GroupCreate(g *models.Group) (*models.Group, error)
	GroupFind(g *models.Group) (*models.Group, error)
	GroupsFind() ([]*models.Group, error)
	GroupDelete(g *models.Group) (*models.Group, error)
	GroupUpdate(g *models.Group) (*models.Group, error)
	GroupDocInsert(g *models.Group) (*models.Group, error)
}
