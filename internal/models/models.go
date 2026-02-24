package models

type PingResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}