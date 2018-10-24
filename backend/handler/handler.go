package handler

import (
	"net/http"
	"strings"

	"github.com/sahandhnj/apiclient/backend/handler/model"
)

type Handler struct {
	Model *model.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/api/model"):
		http.StripPrefix("/api", h.Model).ServeHTTP(w, r)
	}
}
