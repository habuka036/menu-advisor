package service

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/habuka036/menu-advisor/internal/models"
)

// DocumentProcessor handles processing of various document types
type DocumentProcessor struct {
	menuService *MenuAdvisorService
}

// NewDocumentProcessor creates a new document processor
func NewDocumentProcessor(menuService *MenuAdvisorService) *DocumentProcessor {
	return &DocumentProcessor{
		menuService: menuService,
	}
}

// ProcessDocument processes a document and extracts menu information
func (dp *DocumentProcessor) ProcessDocument(req *models.DocumentProcessingRequest) (*models.DocumentSource, error) {
	// Generate unique ID for this document
	docID := generateDocumentID()
	
	// Create document source record
	doc := &models.DocumentSource{
		ID:           docID,
		Type:         req.Type,
		OriginalName: req.Header.Filename,
		UploadedAt:   time.Now(),
		Status:       "processing",
	}

	// Determine document type if not specified
	if req.Type == "" {
		detectedType, err := dp.detectDocumentType(req.Header.Filename)
		if err != nil {
			doc.Status = "error"
			doc.ErrorMessage = fmt.Sprintf("Failed to detect document type: %v", err)
			return doc, err
		}
		doc.Type = detectedType
	}

	// Process based on document type
	extractedData, err := dp.extractDataFromDocument(req, doc)
	if err != nil {
		doc.Status = "error"
		doc.ErrorMessage = fmt.Sprintf("Failed to extract data: %v", err)
		return doc, err
	}

	// Parse extracted data into menu structure
	menus, err := dp.parseExtractedMenuData(extractedData)
	if err != nil {
		doc.Status = "error"
		doc.ErrorMessage = fmt.Sprintf("Failed to parse menu data: %v", err)
		return doc, err
	}

	// Add parsed menus to the service
	for _, menu := range menus {
		dp.menuService.AddSchoolLunchMenu(menu)
	}

	// Mark as completed
	now := time.Now()
	doc.ProcessedAt = &now
	doc.Status = "completed"

	return doc, nil
}

// detectDocumentType determines the document type based on file extension
func (dp *DocumentProcessor) detectDocumentType(filename string) (models.DocumentType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".json":
		return models.DocumentTypeJSON, nil
	case ".pdf":
		// For now, default to PDF text - could enhance to detect if image-based
		return models.DocumentTypePDFText, nil
	case ".jpg", ".jpeg", ".png", ".bmp", ".gif":
		return models.DocumentTypeImage, nil
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// extractDataFromDocument extracts raw data based on document type
func (dp *DocumentProcessor) extractDataFromDocument(req *models.DocumentProcessingRequest, doc *models.DocumentSource) (*models.ExtractedMenuData, error) {
	switch doc.Type {
	case models.DocumentTypeJSON:
		return dp.extractFromJSON(req.File, doc.ID)
	case models.DocumentTypePDFText:
		return dp.extractFromPDFText(req.File, doc.ID)
	case models.DocumentTypePDFImage:
		return dp.extractFromPDFImage(req.File, doc.ID)
	case models.DocumentTypeImage:
		return dp.extractFromImage(req.File, doc.ID)
	default:
		return nil, fmt.Errorf("unsupported document type: %s", doc.Type)
	}
}

// extractFromJSON handles JSON file processing (existing functionality)
func (dp *DocumentProcessor) extractFromJSON(file multipart.File, sourceID string) (*models.ExtractedMenuData, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	return &models.ExtractedMenuData{
		SourceID:    sourceID,
		RawText:     string(data),
		ExtractedAt: time.Now(),
		Confidence:  1.0, // JSON is always 100% confident
		Metadata:    map[string]string{"format": "json"},
	}, nil
}

// extractFromPDFText extracts text from text-based PDFs
func (dp *DocumentProcessor) extractFromPDFText(file multipart.File, sourceID string) (*models.ExtractedMenuData, error) {
	// Placeholder implementation - would use a PDF library like unidoc/unipdf
	// For now, return a mock implementation
	return &models.ExtractedMenuData{
		SourceID:    sourceID,
		RawText:     "PDF text extraction not yet implemented",
		ExtractedAt: time.Now(),
		Confidence:  0.0,
		Metadata:    map[string]string{"format": "pdf_text", "status": "not_implemented"},
	}, fmt.Errorf("PDF text extraction not yet implemented")
}

// extractFromPDFImage extracts text from image-based PDFs using OCR
func (dp *DocumentProcessor) extractFromPDFImage(file multipart.File, sourceID string) (*models.ExtractedMenuData, error) {
	// Placeholder implementation - would use OCR library like tesseract
	return &models.ExtractedMenuData{
		SourceID:    sourceID,
		RawText:     "PDF image OCR not yet implemented",
		ExtractedAt: time.Now(),
		Confidence:  0.0,
		Metadata:    map[string]string{"format": "pdf_image", "status": "not_implemented"},
	}, fmt.Errorf("PDF image OCR not yet implemented")
}

// extractFromImage extracts text from image files using OCR
func (dp *DocumentProcessor) extractFromImage(file multipart.File, sourceID string) (*models.ExtractedMenuData, error) {
	// Placeholder implementation - would use OCR library like tesseract
	return &models.ExtractedMenuData{
		SourceID:    sourceID,
		RawText:     "Image OCR not yet implemented",
		ExtractedAt: time.Now(),
		Confidence:  0.0,
		Metadata:    map[string]string{"format": "image", "status": "not_implemented"},
	}, fmt.Errorf("Image OCR not yet implemented")
}

// parseExtractedMenuData converts extracted raw data into structured menu data
func (dp *DocumentProcessor) parseExtractedMenuData(data *models.ExtractedMenuData) ([]models.SchoolLunchMenu, error) {
	// For JSON format, use existing parsing logic
	if data.Metadata["format"] == "json" {
		return dp.parseJSONMenuData(data.RawText)
	}
	
	// For other formats, would implement text parsing logic here
	// This would involve NLP/pattern matching to extract menu information
	// from OCR text or PDF text
	
	return nil, fmt.Errorf("parsing for format %s not yet implemented", data.Metadata["format"])
}

// parseJSONMenuData parses JSON menu data (reuses existing logic)
func (dp *DocumentProcessor) parseJSONMenuData(jsonText string) ([]models.SchoolLunchMenu, error) {
	var menus []models.SchoolLunchMenu
	if err := json.Unmarshal([]byte(jsonText), &menus); err != nil {
		return nil, fmt.Errorf("failed to parse JSON menu data: %w", err)
	}
	return menus, nil
}

// generateDocumentID generates a unique ID for a document
func generateDocumentID() string {
	return fmt.Sprintf("doc_%d", time.Now().UnixNano())
}