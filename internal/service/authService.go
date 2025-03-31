package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abhayishere/lokrr-api-gateway/internal/models"
	"github.com/abhayishere/lokrr-proto/gen/authpb"
)

func (h *serviceImpl) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req authpb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:        400,
			Description: "Invalid request body",
		})
		return
	}

	res, err := h.AuthClient.RegisterUser(context.Background(), &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:        500,
			Description: fmt.Sprintf("Error in registering user: %v", err),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *serviceImpl) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req authpb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:        400,
			Description: "Invalid request body",
		})
		return
	}

	res, err := h.AuthClient.LoginUser(context.Background(), &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // Use 401 for authentication errors
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:        401,
			Description: fmt.Sprintf("Error in logging in user: %v", err),
		})
		return
	}
	json.NewEncoder(w).Encode(res)
}
