package model

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/sahandhnj/apiclient/backend/middleware"
	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
	DbHandler *db.DBStore
}

func NewHandler(dbHandler *db.DBStore) *Handler {
	h := &Handler{
		Router:    mux.NewRouter(),
		DbHandler: dbHandler,
	}

	fmt.Println("Setting up Model routes")
	h.Handle("/model/test", middleware.Chain(h.helloWorldHandler, middleware.Logging()))
	h.Handle("/model/{modelname}/version/{versionNumber}", middleware.Chain(h.proxyToApi, middleware.Logging()))

	return h
}

func (handler *Handler) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func (handler *Handler) proxyToApi(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	modelName := vars["modelname"]
	versionNumber, err := strconv.Atoi(vars["versionNumber"])

	if err != nil {
		log.Fatal(err)
	}

	modelservice, err := service.NewModelService(nil, handler.DbHandler)
	if err != nil {
		log.Fatal(err)
	}

	model, err := modelservice.DBHandler.ModelService.ModelByName(modelName)
	if err != nil {
		log.Fatal(err)
	}

	versionService, err := service.NewVersionService(model, handler.DbHandler)
	if err != nil {
		log.Fatal(err)
	}

	version, err := versionService.DBHandler.VersionService.VersionByVersionNumber(versionNumber, model.ID)

	url := "http://127.0.0.1:" + strconv.Itoa(version.Port)
	serveReverseProxy(url, res, req)
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Path = "/predict"
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// req.Host = url.Host
	// req.RequestURI = "/t"

	// fmt.Println(url)
	// fmt.Println(url.Host)
	// fmt.Println(url.Scheme)
	// fmt.Println(req.Header.Get("Host"))
	// fmt.Println("----------")
	// fmt.Printf("%+v\n", req)

	proxy.ServeHTTP(res, req)
}
