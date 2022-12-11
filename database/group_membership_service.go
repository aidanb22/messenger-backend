package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/ablancas22/messenger-backend/models"
	"time"
)

// GroupService is used by the app to manage all group related controllers and functionality
type GroupMembershipService struct {
	collection DBCollection
	db         DBClient
	handler    *DBHandler[*groupMembershipModel]
}

// NewGroupService is an exported function used to initialize a new GroupService struct
func NewGroupMembershipService(db DBClient, handler *DBHandler[*groupMembershipModel]) *GroupMembershipService {
	collection := db.GetCollection("group_memberships")
	return &GroupMembershipService{collection, db, handler}
}

// GroupCreate is used to create a new user group
func (p *GroupMembershipService) GroupMembershipCreate(g *models.GroupMembership) (*models.GroupMembership, error) {
	err := g.Validate("create")
	if err != nil {
		return nil, err
	}
	gm, err := newGroupMembershipModel(g)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\npreGMID", gm.GroupId, gm.UserId)
	gRes, err := p.handler.FindOne(&groupMembershipModel{GroupId: gm.GroupId, UserId: gm.UserId})
	fmt.Println("\n\npostGMID", gRes, err)

	if err == nil {
		return nil, errors.New("group name exists")
	}
	gm, err = p.handler.InsertOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupsFind is used to find all group docs in a MongoDB Collection
func (p *GroupMembershipService) GroupMembershipsFind(g *models.GroupMembership) ([]*models.GroupMembership, error) {
	var groups []*models.GroupMembership
	//Todo: build filter using g
	gms, err := p.handler.FindMany(&groupMembershipModel{})
	if err != nil {
		return groups, err
	}
	for _, gm := range gms {
		groups = append(groups, gm.toRoot())
	}
	return groups, nil
}

// GroupFind is used to find a specific group doc
func (p *GroupMembershipService) GroupMembershipFind(g *models.GroupMembership) (*models.GroupMembership, error) {
	gm, err := newGroupMembershipModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.handler.FindOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupDelete is used to delete a group doc
func (p *GroupMembershipService) GroupMembershipDelete(g *models.GroupMembership) (*models.GroupMembership, error) {
	gm, err := newGroupMembershipModel(g)
	if err != nil {
		return nil, err
	}
	gm, err = p.handler.DeleteOne(gm)
	if err != nil {
		return nil, err
	}
	return gm.toRoot(), err
}

// GroupUpdate is used to update an existing group
func (p *GroupMembershipService) GroupMembershipUpdate(g *models.GroupMembership) (*models.GroupMembership, error) {
	var filter models.GroupMembership
	err := g.Validate("create")
	if err != nil {
		return nil, errors.New("missing valid query filter")
	}
	filter.Id = g.Id
	f, err := newGroupMembershipModel(&filter)
	if err != nil {
		return nil, err
	}
	gm, err := newGroupMembershipModel(g)
	if err != nil {
		return nil, err
	}
	_, groupErr := p.handler.FindOne(f)
	if groupErr != nil {
		return nil, errors.New("group not found")
	}
	gm, err = p.handler.UpdateOne(f, gm)
	return gm.toRoot(), err
}

// GroupDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (p *GroupMembershipService) GroupMembershipDocInsert(g *models.GroupMembership) (*models.GroupMembership, error) {
	insertGroup, err := newGroupMembershipModel(g)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = p.collection.InsertOne(ctx, insertGroup)
	if err != nil {
		return nil, err
	}
	return insertGroup.toRoot(), nil
}
