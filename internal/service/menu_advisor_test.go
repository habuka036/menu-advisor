package service

import (
	"testing"
	"time"

	"github.com/habuka036/menu-advisor/internal/models"
)

func TestNewMenuAdvisorService(t *testing.T) {
	service := NewMenuAdvisorService()
	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
	if service.homeMenuDB == nil {
		t.Error("Expected homeMenuDB to be initialized")
	}
}

func TestGenerateHomeMenuSuggestion(t *testing.T) {
	service := NewMenuAdvisorService()
	
	// Add test school lunch data
	testLunch := models.SchoolLunchMenu{
		Date:       time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
		MainDish:   "鶏肉の照り焼き",
		SideDishes: []string{"野菜炒め", "白米"},
		Soup:       "味噌汁（わかめ）",
		Nutrition: models.Nutrition{
			Calories: 650,
			Protein:  28.5,
		},
	}
	service.schoolLunches = []models.SchoolLunchMenu{testLunch}

	// Test breakfast suggestion
	breakfast, err := service.GenerateHomeMenuSuggestion(testLunch.Date, "breakfast")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if breakfast.MealType != "breakfast" {
		t.Errorf("Expected meal type 'breakfast', got: %s", breakfast.MealType)
	}
	if breakfast.MainDish == "" {
		t.Error("Expected main dish to be set")
	}
	if len(breakfast.SideDishes) == 0 {
		t.Error("Expected side dishes to be set")
	}
	if breakfast.Reason == "" {
		t.Error("Expected reason to be provided")
	}

	// Test dinner suggestion
	dinner, err := service.GenerateHomeMenuSuggestion(testLunch.Date, "dinner")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if dinner.MealType != "dinner" {
		t.Errorf("Expected meal type 'dinner', got: %s", dinner.MealType)
	}
	if dinner.MainDish == "" {
		t.Error("Expected main dish to be set")
	}
}

func TestGetSchoolLunchForDate(t *testing.T) {
	service := NewMenuAdvisorService()
	
	testDate := time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)
	testLunch := models.SchoolLunchMenu{
		Date:     testDate,
		MainDish: "鶏肉の照り焼き",
	}
	service.schoolLunches = []models.SchoolLunchMenu{testLunch}

	// Test existing date
	lunch, err := service.GetSchoolLunchForDate(testDate)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if lunch.MainDish != "鶏肉の照り焼き" {
		t.Errorf("Expected main dish '鶏肉の照り焼き', got: %s", lunch.MainDish)
	}

	// Test non-existing date
	wrongDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
	_, err = service.GetSchoolLunchForDate(wrongDate)
	if err == nil {
		t.Error("Expected error for non-existing date")
	}
}