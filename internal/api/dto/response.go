package dto

// APIResponse is the standard response structure for both success and error responses.
type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Timestamp  string      `json:"timestamp,omitempty"`
}
