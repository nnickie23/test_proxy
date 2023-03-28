package handler

import (
	"net/http"

	_ "github.com/nnickie23/test_proxy/docs"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

// InitRoutes - initialize API routes
func (a *handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	// Swagger documentation
	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/task", a.hTaskCreate).Methods("POST")
	r.HandleFunc("/task/{taskID}", a.hTaskGetStatus).Methods("GET")

	return a.middleware(r)
}
