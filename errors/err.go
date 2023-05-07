package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError interface {
	// APIError returns an HTTP status code and an API-safe error message.
	APIError() (int, string)
}

type SentinelAPIError struct {
	Message string
	Code    int
	Err     error
}

func (err SentinelAPIError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err SentinelAPIError) Unwrap() error {
	return err.Err // Returns inner error
}

func (err SentinelAPIError) APIError() (int, string) {
	return err.Code, err.Message
}

// Returns the inner most CustomErrorWrapper
func (err SentinelAPIError) Dig() SentinelAPIError {
	var ew SentinelAPIError
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}

func NewErrorWrapper(code int, err error, message string) error {
	return SentinelAPIError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func JSONHandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var apiErr APIError
	if errors.As(err, &apiErr) {
		status, msg := apiErr.APIError()
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal error"})
	}
}
