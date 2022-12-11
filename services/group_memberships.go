package services

import "github.com/ablancas22/messenger-backend/models"

type GroupMembershipService interface {
	GroupMembershipCreate(g *models.GroupMembership) (*models.GroupMembership, error)
	GroupMembershipFind(g *models.GroupMembership) (*models.GroupMembership, error)
	GroupMembershipsFind(g *models.GroupMembership) ([]*models.GroupMembership, error)
	GroupMembershipDelete(g *models.GroupMembership) (*models.GroupMembership, error)
	GroupMembershipUpdate(g *models.GroupMembership) (*models.GroupMembership, error)
	GroupMembershipDocInsert(g *models.GroupMembership) (*models.GroupMembership, error)
}
