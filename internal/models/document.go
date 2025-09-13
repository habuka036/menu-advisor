package models

import (
	"mime/multipart"
	"time"
)

// DocumentType represents the type of document containing menu information
type DocumentType string

const (
	DocumentTypeJSON       DocumentType = "json"
	DocumentTypePDFText    DocumentType = "pdf_text"    // Text-extractable PDF
	DocumentTypePDFImage   DocumentType = "pdf_image"   // Image-based PDF requiring OCR
	DocumentTypeImage      DocumentType = "image"       // Photo/image file requiring OCR
)

// DocumentSource represents a document containing school lunch menu information
type DocumentSource struct {
	ID           string       `json:"id"`
	Type         DocumentType `json:"type"`
	OriginalName string       `json:"original_name"`
	FilePath     string       `json:"file_path,omitempty"`
	UploadedAt   time.Time    `json:"uploaded_at"`
	ProcessedAt  *time.Time   `json:"processed_at,omitempty"`
	Status       string       `json:"status"` // pending, processing, completed, error
	ErrorMessage string       `json:"error_message,omitempty"`
}

// DocumentProcessingRequest represents a request to process a document
type DocumentProcessingRequest struct {
	File     multipart.File        `json:"-"`
	Header   *multipart.FileHeader `json:"-"`
	Type     DocumentType          `json:"type"`
	DateFrom *time.Time            `json:"date_from,omitempty"`
	DateTo   *time.Time            `json:"date_to,omitempty"`
}

// ExtractedMenuData represents the raw extracted data from a document before parsing
type ExtractedMenuData struct {
	SourceID     string            `json:"source_id"`
	RawText      string            `json:"raw_text"`
	ExtractedAt  time.Time         `json:"extracted_at"`
	Confidence   float64           `json:"confidence,omitempty"` // For OCR results
	Metadata     map[string]string `json:"metadata,omitempty"`
}