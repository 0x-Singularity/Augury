package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/0x-Singularity/Augury/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	router := mux.NewRouter()

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

	// Register FAKEula API routes
	routes.SetupRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
