package api

// BaseError defines model for BaseError.
type BaseError struct {
	Code    *string `json:"code,omitempty"`
	Debug   *string `json:"debug,omitempty"`
	Message *string `json:"message,omitempty"`
}

// BaseResponse defines model for BaseResponse.
type BaseResponse struct {
	Data  *map[string]interface{} `json:"data"`
	Error *BaseError              `json:"error,omitempty"`
}
