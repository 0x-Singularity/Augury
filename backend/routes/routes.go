package routes

import (
	"github.com/0x-Singularity/Augury/controllers"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes all API routes
func SetupRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	// New FAKEula API route
	apiRouter.HandleFunc("/ioc/fakeula", controllers.QueryFakeula).Methods("GET")
}
