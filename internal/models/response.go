package models

// JSONResponse is a generic structure for API responses
type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PingResponse represents the response structure for the /ping endpoint.
type PingResponse struct {
	Message string `json:"message"`
}
