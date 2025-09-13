package models

import "time"

// SchoolLunchMenu represents a school lunch menu for a specific day
type SchoolLunchMenu struct {
	Date        time.Time `json:"date"`
	MainDish    string    `json:"main_dish"`
	SideDishes  []string  `json:"side_dishes"`
	Soup        string    `json:"soup,omitempty"`
	Dessert     string    `json:"dessert,omitempty"`
	Nutrition   Nutrition `json:"nutrition"`
}

// HomeMenuSuggestion represents a suggested home menu
type HomeMenuSuggestion struct {
	Date           time.Time `json:"date"`
	MealType       string    `json:"meal_type"` // breakfast, dinner
	MainDish       string    `json:"main_dish"`
	SideDishes     []string  `json:"side_dishes"`
	Soup           string    `json:"soup,omitempty"`
	Reason         string    `json:"reason"`
	SchoolLunchRef string    `json:"school_lunch_ref"`
}

// Nutrition represents nutritional information
type Nutrition struct {
	Calories     int     `json:"calories"`
	Protein      float64 `json:"protein_g"`
	Carbs        float64 `json:"carbs_g"`
	Fat          float64 `json:"fat_g"`
	Fiber        float64 `json:"fiber_g"`
	Sodium       float64 `json:"sodium_mg"`
	Vegetables   int     `json:"vegetables_servings"`
}

// FoodCategory represents different food categories for balancing
type FoodCategory string

const (
	CategoryProtein    FoodCategory = "protein"
	CategoryVegetables FoodCategory = "vegetables"
	CategoryGrains     FoodCategory = "grains"
	CategoryDairy      FoodCategory = "dairy"
	CategoryFruits     FoodCategory = "fruits"
)

// FoodItem represents a food item with its category and nutritional info
type FoodItem struct {
	Name        string       `json:"name"`
	Category    FoodCategory `json:"category"`
	Nutrition   Nutrition    `json:"nutrition"`
	Season      []string     `json:"season,omitempty"`
	Japanese    bool         `json:"japanese"`
}