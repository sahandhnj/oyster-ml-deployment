package model

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sahandhnj/apiclient/backend/middleware"

	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
}

func NewHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	fmt.Println("Setting up Model routes")
	h.Handle("/model", middleware.Chain(h.helloWorldHandler, middleware.Logging()))

	return h
}

func (handler *Handler) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
