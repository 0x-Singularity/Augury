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

	apiRouter.HandleFunc("/ioc/extract", controllers.ExtractFromText).Methods("POST", "OPTIONS")

	apiRouter.HandleFunc("/ioc/oil", controllers.QueryAllOIL).Methods("GET")
	apiRouter.HandleFunc("/ioc/pdns", controllers.QueryPDNS).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/ioc/ldap", controllers.QueryLDAP).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/ioc/geo", controllers.QueryGeoIP).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/ioc/binary", controllers.QueryBinary).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/ioc/vpn", controllers.QueryVPN).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/ioc/cbr", controllers.QueryCBR).Methods("GET", "OPTIONS")
}
