package server

import (
	"github.com/ablancas22/messenger-backend/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Server is a struct that stores the API Apps high level attributes such as the router, config, and services
type Server struct {
	Router                  *mux.Router
	TokenService            *services.TokenService
	UserService             services.UserService
	GroupService            services.GroupService
	MessageService          services.MessageService
	GroupMembershipsService services.GroupMembershipService
	ConversationService     services.ConversationService
	ContactService          services.ContactService
}

// NewServer is a function used to initialize a new Server struct
func NewServer(u services.UserService, g services.GroupService, tt services.MessageService, t *services.TokenService, gm services.GroupMembershipService, c services.ConversationService, co services.ContactService) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router = NewGroupRouter(router, t, g, u, gm)
	router = NewUserRouter(router, t, u, g)
	router = NewMessageRouter(router, t, tt)
	router = NewConversationRouter(router, t, tt, c, u)
	router = NewContactRouter(router, t, tt, co, u)
	return &Server{
		Router:                  router,
		TokenService:            t,
		UserService:             u,
		GroupService:            g,
		MessageService:          tt,
		GroupMembershipsService: gm,
		ConversationService:     c,
		ContactService:          co,
	}
}

// Start starts the initialized Server
func (s *Server) Start() {
	log.Println("Listening on port " + os.Getenv("PORT"))
	go func() {
		if err := http.ListenAndServe(":"+os.Getenv("PORT"), handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)
}
