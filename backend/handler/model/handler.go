package model

import (
	"encoding/json"
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
	DbHandler      *db.DBStore
	VersionService *service.VersionService
	ModelService   *service.ModelService
}

func NewHandler(dbHandler *db.DBStore, vs *service.VersionService, ms *service.ModelService) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		DbHandler:      dbHandler,
		VersionService: vs,
		ModelService:   ms,
	}

	fmt.Println("Setting up Model routes")
	h.Handle("/model/test", middleware.Chain(h.helloWorldHandler, middleware.Logging())).Methods("GET")
	h.Handle("/model/{modelname}/version/{versionNumber}", middleware.Chain(h.proxyToApi, middleware.Logging()))
	h.Handle("/model/all", middleware.Chain(h.getAllModels, middleware.Logging())).Methods("GET")

	return h
}

func (handler *Handler) getAllModels(w http.ResponseWriter, r *http.Request) {
	models, err := handler.ModelService.GetAll()
	if err != nil {
		fmt.Println(err)
	}

	json, err := json.Marshal(models)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (handler *Handler) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func (handler *Handler) proxyToApi(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	modelName := vars["modelname"]
	versionNumber, err := strconv.Atoi(vars["versionNumber"])

	model, err := handler.ModelService.DBHandler.ModelService.ModelByName(modelName)
	if err != nil {
		log.Fatal(err)
	}

	version, err := handler.VersionService.DBHandler.VersionService.VersionByVersionNumber(versionNumber, model.ID)

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
