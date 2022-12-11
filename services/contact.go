package services

import "github.com/ablancas22/messenger-backend/models"

type ContactService interface {
	ContactCreate(g *models.Contact) (*models.Contact, error)
	ContactFind(g *models.Contact) (*models.Contact, error)
	ContactsFind(g *models.Contact) ([]*models.Contact, error)
	ContactDelete(g *models.Contact) (*models.Contact, error)
	ContactDocInsert(g *models.Contact) (*models.Contact, error)
}
