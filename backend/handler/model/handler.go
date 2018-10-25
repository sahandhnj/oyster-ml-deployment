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
	url := "http://127.0.0.1:5001"

	serveReverseProxy(url, res, req)
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Path = "/predict"
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// req.Host = url.Host
	req.RequestURI = "/t"

	fmt.Println(url)
	fmt.Println(url.Host)
	fmt.Println(url.Scheme)
	fmt.Println(req.Header.Get("Host"))
	fmt.Println("----------")
	fmt.Printf("%+v\n", req)

	proxy.ServeHTTP(res, req)
}
