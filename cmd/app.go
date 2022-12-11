package cmd

import (
	"context"
	"github.com/ablancas22/messenger-backend/database"
	"github.com/ablancas22/messenger-backend/models"
	"github.com/ablancas22/messenger-backend/server"
	"github.com/ablancas22/messenger-backend/services"
	"github.com/ablancas22/messenger-backend/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"
)

// App is the highest level struct of the rest_api application. Stores the server, client, and config settings.
type App struct {
	server *server.Server
	db     database.DBClient
}

// Initialize is a function used to initialize a new instantiation of the API Application
func (a *App) Initialize() error {
	var err error
	// 1) Initialize config settings & set environmental variables
	if os.Getenv("ENV") != "docker-dev" {
		conf, err := getConfigurations()
		if err != nil {
			return err
		}
		conf.InitializeEnvironmentalVars()
	}
	// 2) Initialize & Connect DB Client
	a.db, err = database.InitializeNewClient()
	if err != nil {
		return err
	}
	err = a.db.Connect()
	if err != nil {
		return err
	}
	// 3) Initial DB Services
	gHandler := a.db.NewGroupHandler()
	uHandler := a.db.NewUserHandler()
	blHandler := a.db.NewBlacklistHandler()
	tHandler := a.db.NewMessageHandler()
	gmHandler := a.db.NewGroupMembershipHandler()
	cHandler := a.db.NewConversationHandler()
	coHandler := a.db.NewContactHandler()

	gService := database.NewGroupService(a.db, gHandler)
	uService := database.NewUserService(a.db, uHandler, gHandler)
	bService := database.NewBlacklistService(a.db, blHandler)
	gmService := database.NewGroupMembershipService(a.db, gmHandler)
	tService := services.NewTokenService(uService, gService, bService)
	ttService := database.NewMessageService(a.db, tHandler, uHandler, gHandler)
	cService := database.NewConversationService(a.db, cHandler)
	coService := database.NewContactService(a.db, coHandler)

	// 4) Create RootAdmin user if database is empty
	var group models.Group
	var adminUser models.User
	group.Name = os.Getenv("ROOT_GROUP")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := a.db.GetCollection("groups").CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if docCount == 0 {
		group.Id = utilities.GenerateObjectID()
		if err != nil {
			return err
		}
		adminUser.Username = os.Getenv("ROOT_ADMIN")
		adminUser.Email = os.Getenv("ROOT_EMAIL")
		adminUser.Password = os.Getenv("ROOT_PASSWORD")
		_, err = uService.UserCreate(&adminUser)
		if err != nil {
			return err
		}
	}
	// 5) Initialize Server
	a.server = server.NewServer(uService, gService, ttService, tService, gmService, cService, coService)
	return nil
}

// Run is a function used to run a previously initialized API Application
func (a *App) Run() {
	defer a.db.Close()
	a.server.Start()
}
