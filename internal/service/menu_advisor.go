package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/habuka036/menu-advisor/internal/models"
)

// MenuAdvisorService provides menu recommendation functionality
type MenuAdvisorService struct {
	schoolLunches []models.SchoolLunchMenu
	homeMenuDB    map[string][]models.FoodItem
}

// NewMenuAdvisorService creates a new instance of the service
func NewMenuAdvisorService() *MenuAdvisorService {
	service := &MenuAdvisorService{
		homeMenuDB: make(map[string][]models.FoodItem),
	}
	service.initializeHomeMenuDatabase()
	return service
}

// LoadSchoolLunchData loads school lunch data from a JSON file
func (s *MenuAdvisorService) LoadSchoolLunchData(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, &s.schoolLunches); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// GetSchoolLunchForDate returns the school lunch menu for a specific date
func (s *MenuAdvisorService) GetSchoolLunchForDate(date time.Time) (*models.SchoolLunchMenu, error) {
	dateStr := date.Format("2006-01-02")
	for _, lunch := range s.schoolLunches {
		if lunch.Date.Format("2006-01-02") == dateStr {
			return &lunch, nil
		}
	}
	return nil, fmt.Errorf("no school lunch found for date: %s", dateStr)
}

// GenerateHomeMenuSuggestion generates home menu suggestions based on school lunch
func (s *MenuAdvisorService) GenerateHomeMenuSuggestion(date time.Time, mealType string) (*models.HomeMenuSuggestion, error) {
	schoolLunch, err := s.GetSchoolLunchForDate(date)
	if err != nil {
		return nil, err
	}

	suggestion := &models.HomeMenuSuggestion{
		Date:           date,
		MealType:       mealType,
		SchoolLunchRef: schoolLunch.MainDish,
	}

	// Generate complementary menu based on school lunch
	if mealType == "breakfast" {
		s.generateBreakfastSuggestion(suggestion, schoolLunch)
	} else if mealType == "dinner" {
		s.generateDinnerSuggestion(suggestion, schoolLunch)
	}

	return suggestion, nil
}

func (s *MenuAdvisorService) generateBreakfastSuggestion(suggestion *models.HomeMenuSuggestion, schoolLunch *models.SchoolLunchMenu) {
	// For breakfast, focus on lighter options that complement lunch
	if strings.Contains(schoolLunch.MainDish, "魚") {
		suggestion.MainDish = "卵焼き"
		suggestion.SideDishes = []string{"のり", "みそ汁"}
		suggestion.Reason = "昼食で魚を摂取するため、朝食ではタンパク質として卵を提案"
	} else if strings.Contains(schoolLunch.MainDish, "肉") {
		suggestion.MainDish = "焼き魚（アジ）"
		suggestion.SideDishes = []string{"野菜サラダ", "みそ汁"}
		suggestion.Reason = "昼食で肉類を摂取するため、朝食では魚でバランスを取る"
	} else if strings.Contains(schoolLunch.MainDish, "カレー") {
		suggestion.MainDish = "納豆"
		suggestion.SideDishes = []string{"野菜炒め", "みそ汁"}
		suggestion.Reason = "昼食が重めのカレーのため、朝食は軽めの和食で消化を助ける"
	} else {
		suggestion.MainDish = "焼き鮭"
		suggestion.SideDishes = []string{"おひたし", "みそ汁"}
		suggestion.Reason = "栄養バランスを考慮した和食中心の朝食"
	}
}

func (s *MenuAdvisorService) generateDinnerSuggestion(suggestion *models.HomeMenuSuggestion, schoolLunch *models.SchoolLunchMenu) {
	// For dinner, complement what was missing in lunch or provide variety
	if strings.Contains(schoolLunch.MainDish, "鶏") {
		suggestion.MainDish = "魚の煮付け"
		suggestion.SideDishes = []string{"野菜の天ぷら", "白米"}
		suggestion.Soup = "すまし汁"
		suggestion.Reason = "昼食で鶏肉を摂取したため、夕食では魚でタンパク質の種類を変える"
	} else if strings.Contains(schoolLunch.MainDish, "魚") {
		suggestion.MainDish = "豚しゃぶしゃぶ"
		suggestion.SideDishes = []string{"温野菜", "白米"}
		suggestion.Soup = "みそ汁"
		suggestion.Reason = "昼食で魚を摂取したため、夕食では豚肉でタンパク質の種類を変える"
	} else if strings.Contains(schoolLunch.MainDish, "豚") {
		suggestion.MainDish = "鯖の塩焼き"
		suggestion.SideDishes = []string{"筑前煮", "白米"}
		suggestion.Soup = "わかめスープ"
		suggestion.Reason = "昼食で豚肉を摂取したため、夕食では魚でバランスを取る"
	} else if strings.Contains(schoolLunch.MainDish, "カレー") {
		suggestion.MainDish = "鶏の唐揚げ"
		suggestion.SideDishes = []string{"キャベツサラダ", "白米"}
		suggestion.Soup = "みそ汁"
		suggestion.Reason = "昼食がスパイシーなカレーのため、夕食は優しい味付けの料理で胃を休める"
	} else {
		suggestion.MainDish = "牛肉炒め"
		suggestion.SideDishes = []string{"もやし炒め", "白米"}
		suggestion.Soup = "中華スープ"
		suggestion.Reason = "栄養バランスを考慮したボリュームのある夕食"
	}
}

func (s *MenuAdvisorService) initializeHomeMenuDatabase() {
	// Initialize home menu database with common Japanese dishes
	proteins := []models.FoodItem{
		{Name: "焼き鮭", Category: models.CategoryProtein, Japanese: true},
		{Name: "卵焼き", Category: models.CategoryProtein, Japanese: true},
		{Name: "納豆", Category: models.CategoryProtein, Japanese: true},
		{Name: "豚しゃぶしゃぶ", Category: models.CategoryProtein, Japanese: true},
		{Name: "鶏の唐揚げ", Category: models.CategoryProtein, Japanese: true},
		{Name: "魚の煮付け", Category: models.CategoryProtein, Japanese: true},
		{Name: "鯖の塩焼き", Category: models.CategoryProtein, Japanese: true},
		{Name: "牛肉炒め", Category: models.CategoryProtein, Japanese: true},
	}

	vegetables := []models.FoodItem{
		{Name: "野菜サラダ", Category: models.CategoryVegetables, Japanese: false},
		{Name: "おひたし", Category: models.CategoryVegetables, Japanese: true},
		{Name: "野菜炒め", Category: models.CategoryVegetables, Japanese: true},
		{Name: "野菜の天ぷら", Category: models.CategoryVegetables, Japanese: true},
		{Name: "温野菜", Category: models.CategoryVegetables, Japanese: true},
		{Name: "筑前煮", Category: models.CategoryVegetables, Japanese: true},
		{Name: "キャベツサラダ", Category: models.CategoryVegetables, Japanese: false},
		{Name: "もやし炒め", Category: models.CategoryVegetables, Japanese: true},
	}

	grains := []models.FoodItem{
		{Name: "白米", Category: models.CategoryGrains, Japanese: true},
		{Name: "玄米", Category: models.CategoryGrains, Japanese: true},
		{Name: "パン", Category: models.CategoryGrains, Japanese: false},
	}

	s.homeMenuDB["protein"] = proteins
	s.homeMenuDB["vegetables"] = vegetables
	s.homeMenuDB["grains"] = grains
}

// GetAllSchoolLunches returns all loaded school lunch menus
func (s *MenuAdvisorService) GetAllSchoolLunches() []models.SchoolLunchMenu {
	return s.schoolLunches
}

// AddSchoolLunchMenu adds a new school lunch menu to the service
func (s *MenuAdvisorService) AddSchoolLunchMenu(menu models.SchoolLunchMenu) {
	// Check if menu for this date already exists and update it
	dateStr := menu.Date.Format("2006-01-02")
	for i, existingMenu := range s.schoolLunches {
		if existingMenu.Date.Format("2006-01-02") == dateStr {
			s.schoolLunches[i] = menu
			return
		}
	}
	
	// Add new menu if not found
	s.schoolLunches = append(s.schoolLunches, menu)
}

// AddSchoolLunchMenus adds multiple school lunch menus to the service
func (s *MenuAdvisorService) AddSchoolLunchMenus(menus []models.SchoolLunchMenu) {
	for _, menu := range menus {
		s.AddSchoolLunchMenu(menu)
	}
}