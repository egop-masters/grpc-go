package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Request represents the expected structure of the incoming request.
type Request struct {
	Message string `json:"message"`
}

// Response represents the structure of the outgoing response.
type Response struct {
	Reply string `json:"reply"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON request
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the response
	res := Response{
		Reply: fmt.Sprintf("Hello %s", req.Message),
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
