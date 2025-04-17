package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/0x-Singularity/Augury/models" // Import database models
	"github.com/0x-Singularity/Augury/routes" // Import API routes
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// CORSMiddleware adds CORS headers to the response
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Name")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Initialize Azure SQL database connection
	err = models.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(CORSMiddleware)

	staticFileDirectory := http.Dir("../frontend/static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Could not load template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Could not render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Register API routes
	routes.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	// Debugging: Print all registered routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Registered Route:", path)
		}
		return nil
	})
}
