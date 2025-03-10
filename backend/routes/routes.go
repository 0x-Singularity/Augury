package routes

import (
	"github.com/0x-Singularity/Augury/controllers"
	"github.com/gorilla/mux"
)

// SetupRoutes registers API endpoints on the provided router
func SetupRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Map API paths to controller functions
	apiRouter.HandleFunc("/ioc/lookup", controllers.LookupIOC).Methods("GET")
	//apiRouter.HandleFunc("/ioc/log", controllers.LogIOC).Methods("POST")
	apiRouter.HandleFunc("/ioc/fakeula", controllers.QueryFakeula).Methods("GET")
}
