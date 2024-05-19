package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Request represents the expected structure of the incoming request.
type Request struct {
	Message string `json:"message"`
}

// Response represents the structure of the outgoing response.
type Response struct {
	Reply string `json:"reply"`
}

// Create a new counter vector for counting HTTP requests
var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "masters_http_requests",
		Help: "Number of HTTP requests received.",
	},
	[]string{"path"},
)

func init() {
	// Register the counter with Prometheus
	prometheus.MustRegister(httpRequests)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Increment the counter
	httpRequests.WithLabelValues(r.URL.Path).Inc()

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

	// Expose the registered metrics at /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
