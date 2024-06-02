package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	pb "golangrpc/proto" // Adjust this import path if necessary

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Response struct represents the HTTP response format
type Response struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func main() {
	// Create gRPC connection
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	c := pb.NewEmployeeServiceClient(conn)

	// HTTP handler function
	http.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method supported", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Parse JSON request body
		var request pb.CreateEmployeeRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
			return
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Invoke gRPC method
		response, err := c.CreateEmployee(ctx, &request)
		if err != nil {
			sendJSONResponse(w, false, http.StatusInternalServerError, "Failed to create employee", nil)
			return
		}

		// Marshal gRPC response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			sendJSONResponse(w, false, http.StatusInternalServerError, "Failed to marshal response to JSON", nil)
			return
		}

		// Send success response with data
		sendJSONResponse(w, true, http.StatusOK, "Created employee successfully", responseJSON)
	})

	// Start HTTP server
	log.Println("Starting HTTP server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// sendJSONResponse writes the provided response data as JSON to the ResponseWriter
func sendJSONResponse(w http.ResponseWriter, success bool, code int, message string, data json.RawMessage) {
	response := Response{
		Success: success,
		Code:    code,
		Message: message,
		Data:    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJSON)
}
