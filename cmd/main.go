package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/habuka036/menu-advisor/internal/service"
	"github.com/habuka036/menu-advisor/internal/web"
)

func main() {
	// Initialize the menu advisor service
	menuService := service.NewMenuAdvisorService()

	// Load sample school lunch data
	dataPath := filepath.Join("data", "school_lunch_sample.json")
	if err := menuService.LoadSchoolLunchData(dataPath); err != nil {
		log.Printf("Warning: Could not load school lunch data: %v", err)
		log.Println("The service will run with no school lunch data")
	} else {
		log.Println("Successfully loaded school lunch data")
	}

	// Create HTTP handler
	handler := web.NewHandler(menuService)

	// Set up routes
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/api/suggest", handler.SuggestHandler)
	http.HandleFunc("/api/school-lunches", handler.SchoolLunchHandler)

	// Serve static files if they exist
	staticDir := "web/static"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🍱 学校給食メニューアドバイザー starting on port %s", port)
	log.Printf("📱 Access the service at: http://localhost:%s", port)
	log.Printf("🔗 API endpoints:")
	log.Printf("   GET / - Main web interface")
	log.Printf("   GET /api/suggest?date=YYYY-MM-DD&meal_type=breakfast|dinner")
	log.Printf("   GET /api/school-lunches - All school lunch data")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}