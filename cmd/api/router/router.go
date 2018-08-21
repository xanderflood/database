package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xanderflood/accountant/lib/tools"
	"github.com/xanderflood/accountant/pkg/middleware"
)

//Routes collection of all routes
//go:generate counterfeiter . Routes
type Routes interface {
	CreateTable(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
}

type Router func(w http.ResponseWriter, r *http.Request)

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router(w, r)
}

//New returns a new router for the accountant api
func New(
	routes Routes,
	log tools.Logger,
	publicAssetsPath string,
) Router {
	r := mux.NewRouter()

	res := routes.Resolver

	//static assets
	r.PathPrefix("/public/").Handler(
		http.StripPrefix(
			"/public/",
			http.FileServer(
				http.Dir(publicAssetsPath),
			),
		),
	)

	//////transactions//////

	//index
	r.HandleFunc(
		res.CreateTable.Route(),
		routes.CreateTable(w, r),
	).Methods("GET")
	r.HandleFunc(
		res.Index.Route(),
		routes.Index(w, r),
	).Methods("GET")
	r.HandleFunc(
		res.Insert.Route(),
		routes.Insert(w, r),
	).Methods("GET")

	//TODO handler for 404s?
	return middleware.Wrap(log, r)
}
