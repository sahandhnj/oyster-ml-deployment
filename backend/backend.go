package backend

import (
	"fmt"
	"net/http"

	"github.com/sahandhnj/apiclient/backend/handler"
	"github.com/sahandhnj/apiclient/backend/handler/model"
	"github.com/sahandhnj/apiclient/db"
)

type Server struct {
	Address   string
	Handler   *handler.Handler
	DbHandler *db.DBStore
}

func (server *Server) Start() error {
	model := model.NewHandler(server.DbHandler)

	server.Handler = &handler.Handler{
		Model: model,
	}

	fmt.Println("Listening on: " + server.Address)
	return http.ListenAndServe(server.Address, server.Handler)
}
