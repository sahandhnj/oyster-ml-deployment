package model

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

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
	h.Handle("/proxy", middleware.Chain(h.handleRequestAndRedirect, middleware.Logging()))

	return h
}

func (handler *Handler) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func (handler *Handler) handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	url := "http://127.0.0.1:5001/predict"

	serveReverseProxy(url, res, req)
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}
