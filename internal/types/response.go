package types

import "obs/internal/models"

type Response struct {
	StatusCode int                      `json:"status_code"`
	Success    bool                     `json:"success"`
	Message    string                   `json:"message,omitempty"`
	Data       map[string]any           `json:"data,omitempty"`
	Error      string                   `json:"error,omitempty"`
	Blogs      map[string][]models.Blog `json:"messages,omitempty"`
}
