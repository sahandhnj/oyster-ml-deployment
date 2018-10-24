package backend

import (
	"fmt"
	"net/http"

	"github.com/sahandhnj/apiclient/backend/handler"
	"github.com/sahandhnj/apiclient/backend/handler/model"
)

type Server struct {
	Address string
	Handler *handler.Handler
}

func (server *Server) Start() error {
	model := model.NewHandler()

	server.Handler = &handler.Handler{
		Model: model,
	}

	fmt.Println("Listening on: " + server.Address)
	return http.ListenAndServe(server.Address, server.Handler)
}
