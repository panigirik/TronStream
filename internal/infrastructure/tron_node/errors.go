package tron_node

import (
	"fmt"
	"net/http"
)

// APIError — ошибка, которую вернул сам TronGrid (в отличие от сетевой ошибки).
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("trongrid: HTTP %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode == http.StatusForbidden
}
