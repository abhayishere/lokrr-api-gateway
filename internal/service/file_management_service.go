package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abhayishere/lokrr-api-gateway/internal/middleware"
	"github.com/abhayishere/lokrr-api-gateway/internal/models"
	fmpb "github.com/abhayishere/lokrr-proto/gen/file_management"
)

func (h *serviceImpl) UploadDocument(w http.ResponseWriter, r *http.Request) (*fmpb.UploadResponse, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method")
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return nil, fmt.Errorf("user ID not found in context")
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in parsing form: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in parsing form: %v", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("error in getting file: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in getting file: %v", err)
	}
	defer file.Close()

	stream, err := h.DocClient.UploadDocument(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("error in uploading document: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in uploading document: %v", err)
	}
	buffer := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, fmt.Sprintf("error in reading file: %v", err), http.StatusInternalServerError)
			return nil, fmt.Errorf("error in reading file: %v", err)
		}
		if n == 0 {
			break
		}

		err = stream.Send(&fmpb.UploadRequest{
			UserId:       userID,
			DocumentName: header.Filename,
			DocumentType: header.Header.Get("Content-Type"),
			Chunk:        buffer[:n],
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("error in sending chunk: %v", err), http.StatusInternalServerError)
			return nil, fmt.Errorf("error in sending chunk: %v", err)
		}
	}

	uploadRes, err := stream.CloseAndRecv()
	if err != nil {
		http.Error(w, fmt.Sprintf("error in uploading document: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in uploading document: %v", err)
	}

	return uploadRes, nil
}

func (h *serviceImpl) GetDocument(w http.ResponseWriter, r *http.Request) (*fmpb.GetDocumentResponse, error) {
	// Extract document ID from request
	// For simplicity, hardcoding value here
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method")
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return nil, fmt.Errorf("user ID not found in context")
	}
	var req models.GetDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid request body: %v", err)
	}
	getDocRes, err := h.DocClient.GetDocument(context.Background(), &fmpb.GetDocumentRequest{
		UserId:     userID,
		DocumentId: req.DocumentID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("error in getting document: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in getting document: %v", err)
	}
	return getDocRes, nil
}

func (h *serviceImpl) ListDocuments(w http.ResponseWriter, r *http.Request) (*fmpb.ListDocumentsResponse, error) {
	// Extract user ID from request
	// For simplicity, hardcoding value here
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method")
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return nil, fmt.Errorf("user ID not found in context")
	}

	listDocsRes, err := h.DocClient.ListDocument(context.Background(), &fmpb.ListDocumentsRequest{
		UserId: userID,
	})
	
	if err != nil {
		http.Error(w, fmt.Sprintf("error in listing documents: %v", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error in listing documents: %v", err)
	}
	return listDocsRes, nil
}

func (h *serviceImpl) DeleteDocument(w http.ResponseWriter, r *http.Request) (*fmpb.DeleteDocumentResponse, error) {
	// Extract document ID from request
	// For simplicity, hardcoding value here
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method")
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return nil, fmt.Errorf("user ID not found in context")
	}
	var req models.DeleteDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	deleteDocRes, err := h.DocClient.DeleteDocument(context.Background(), &fmpb.DeleteDocumentRequest{
		UserId:     userID,
		DocumentId: req.DocumentID,
	})
	if err != nil {
		return nil, fmt.Errorf("error in deleting document: %v", err)
	}
	return deleteDocRes, nil
}
