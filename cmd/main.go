package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/abhayishere/lokrr-api-gateway/internal/config"
	"github.com/abhayishere/lokrr-api-gateway/internal/middleware"
	"github.com/abhayishere/lokrr-api-gateway/internal/models"
	"github.com/abhayishere/lokrr-api-gateway/internal/service"
	"github.com/gorilla/handlers"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create gRPC connections
	authConn, err := grpc.Dial(cfg.Grpc.AuthServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Issue in connection to auth server")
	}
	defer authConn.Close()

	docConn, err := grpc.Dial(cfg.Grpc.DocServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Issue in connection to document server")
	}
	defer docConn.Close()

	srv := service.NewService(authConn, docConn)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", srv.RegisterUser)
	mux.HandleFunc("/login", srv.LoginUser)

	mux.Handle("/upload", middleware.AuthMiddleware(srv.AuthClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := srv.UploadDocument(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Code: 500, Description: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	})))

	mux.Handle("/get", middleware.AuthMiddleware(srv.AuthClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := srv.GetDocument(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Code: 500, Description: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	})))

	mux.Handle("/list", middleware.AuthMiddleware(srv.AuthClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := srv.ListDocuments(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Code: 500, Description: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	})))

	mux.Handle("/delete", middleware.AuthMiddleware(srv.AuthClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := srv.DeleteDocument(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Code: 500, Description: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	})))

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(mux)

	log.Printf("Starting server on :%s", cfg.Server.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.Server.HTTPPort, corsHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

/*
each grpc server has proto file
running proto file will generate handler functions
these handler functions are imported in gateway service
then from gateway service these handlers are called by converted json request to protobuf,
whose models are also defined in proto file of the respective services
*/
