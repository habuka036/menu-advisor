package service

import (
	"mime/multipart"
	"strings"
	"testing"

	"github.com/habuka036/menu-advisor/internal/models"
)

// mockFile implements multipart.File for testing
type mockFile struct {
	*strings.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func newMockFile(data string) multipart.File {
	return &mockFile{strings.NewReader(data)}
}

func TestNewDocumentProcessor(t *testing.T) {
	menuService := NewMenuAdvisorService()
	processor := NewDocumentProcessor(menuService)

	if processor == nil {
		t.Error("Expected processor to be created, got nil")
	}
	if processor.menuService == nil {
		t.Error("Expected menuService to be set")
	}
}

func TestDetectDocumentType(t *testing.T) {
	processor := &DocumentProcessor{}

	tests := []struct {
		filename     string
		expectedType models.DocumentType
		expectError  bool
	}{
		{"menu.json", models.DocumentTypeJSON, false},
		{"menu.pdf", models.DocumentTypePDFText, false},
		{"menu.jpg", models.DocumentTypeImage, false},
		{"menu.jpeg", models.DocumentTypeImage, false},
		{"menu.png", models.DocumentTypeImage, false},
		{"menu.txt", "", true},
		{"menu", "", true},
	}

	for _, test := range tests {
		result, err := processor.detectDocumentType(test.filename)
		
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for filename %s, got none", test.filename)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for filename %s: %v", test.filename, err)
			}
			if result != test.expectedType {
				t.Errorf("Expected type %s for filename %s, got %s", test.expectedType, test.filename, result)
			}
		}
	}
}

func TestExtractFromJSON(t *testing.T) {
	menuService := NewMenuAdvisorService()
	processor := NewDocumentProcessor(menuService)

	jsonData := `[{
		"date": "2025-01-20T00:00:00Z",
		"main_dish": "ハンバーグ",
		"side_dishes": ["野菜サラダ", "白米"],
		"soup": "コンソメスープ",
		"nutrition": {
			"calories": 700,
			"protein_g": 30.0,
			"carbs_g": 80.0,
			"fat_g": 25.0,
			"fiber_g": 5.0,
			"sodium_mg": 900,
			"vegetables_servings": 2
		}
	}]`

	file := newMockFile(jsonData)

	result, err := processor.extractFromJSON(file, "test_id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.SourceID != "test_id" {
		t.Errorf("Expected SourceID 'test_id', got '%s'", result.SourceID)
	}
	if result.RawText != jsonData {
		t.Errorf("Expected RawText to match input JSON")
	}
	if result.Confidence != 1.0 {
		t.Errorf("Expected Confidence 1.0, got %f", result.Confidence)
	}
	if result.Metadata["format"] != "json" {
		t.Errorf("Expected format 'json', got '%s'", result.Metadata["format"])
	}
}

func TestParseJSONMenuData(t *testing.T) {
	menuService := NewMenuAdvisorService()
	processor := NewDocumentProcessor(menuService)

	jsonData := `[{
		"date": "2025-01-20T00:00:00Z",
		"main_dish": "ハンバーグ",
		"side_dishes": ["野菜サラダ", "白米"],
		"soup": "コンソメスープ",
		"nutrition": {
			"calories": 700,
			"protein_g": 30.0,
			"carbs_g": 80.0,
			"fat_g": 25.0,
			"fiber_g": 5.0,
			"sodium_mg": 900,
			"vegetables_servings": 2
		}
	}]`

	menus, err := processor.parseJSONMenuData(jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(menus) != 1 {
		t.Fatalf("Expected 1 menu, got %d", len(menus))
	}

	menu := menus[0]
	if menu.MainDish != "ハンバーグ" {
		t.Errorf("Expected MainDish 'ハンバーグ', got '%s'", menu.MainDish)
	}
	if len(menu.SideDishes) != 2 {
		t.Errorf("Expected 2 side dishes, got %d", len(menu.SideDishes))
	}
	if menu.Soup != "コンソメスープ" {
		t.Errorf("Expected Soup 'コンソメスープ', got '%s'", menu.Soup)
	}
}

func TestProcessDocumentJSON(t *testing.T) {
	menuService := NewMenuAdvisorService()
	processor := NewDocumentProcessor(menuService)

	jsonData := `[{
		"date": "2025-01-20T00:00:00Z",
		"main_dish": "ハンバーグ",
		"side_dishes": ["野菜サラダ", "白米"],
		"soup": "コンソメスープ",
		"nutrition": {
			"calories": 700,
			"protein_g": 30.0,
			"carbs_g": 80.0,
			"fat_g": 25.0,
			"fiber_g": 5.0,
			"sodium_mg": 900,
			"vegetables_servings": 2
		}
	}]`

	// Create file header
	header := &multipart.FileHeader{
		Filename: "test.json",
	}

	// Create request
	req := &models.DocumentProcessingRequest{
		File:   newMockFile(jsonData),
		Header: header,
		Type:   models.DocumentTypeJSON,
	}

	// Test processing
	result, err := processor.ProcessDocument(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}
	if result.Type != models.DocumentTypeJSON {
		t.Errorf("Expected type JSON, got %s", result.Type)
	}
	if result.OriginalName != "test.json" {
		t.Errorf("Expected original name 'test.json', got '%s'", result.OriginalName)
	}

	// Check that menu was added to service
	menus := menuService.GetAllSchoolLunches()
	found := false
	for _, menu := range menus {
		if menu.MainDish == "ハンバーグ" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected new menu to be added to service")
	}
}