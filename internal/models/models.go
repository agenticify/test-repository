package models

// PingResponse represents the response structure for the ping endpoint.
type PingResponse struct {
	Message string `json:"message"`
}

// HealthCheckResponse represents the response structure for health check endpoint.
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
