package models

// PingResponse represents the response structure for the ping endpoint.
type PingResponse struct {
	Message string `json:"message"`
}

// HealthResponse represents the response structure for the health endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
