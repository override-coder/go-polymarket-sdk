package http

import "fmt"

type UpstreamServiceError struct {
	Service    string
	StatusCode int
	Message    string
	Body       any
}

func (e *UpstreamServiceError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("[%s] %d %s", e.Service, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Service, e.Message)
}
