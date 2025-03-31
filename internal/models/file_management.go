package models

type UploadDocumentRequest struct {
	DocumentName string `json:"document_name"`
	DocumentType string `json:"document_type"`
}

type GetDocumentRequest struct {
	DocumentID string `json:"document_id"`
}

type DeleteDocumentRequest struct {
	DocumentID string `json:"document_id"`
}
