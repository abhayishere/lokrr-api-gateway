package models

type ErrorResponse struct {
	Code        int    `json:"code"`        // Error code
	Description string `json:"description"` // Error description
}
