package models

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type PingResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
