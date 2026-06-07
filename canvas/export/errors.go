package export

import "errors"

// Common errors returned by providers.
var (
	// ErrNotFound is returned when a canvas does not exist in the external system.
	ErrNotFound = errors.New("canvas not found")

	// ErrUnsupportedType is returned when a provider doesn't support the canvas type.
	ErrUnsupportedType = errors.New("unsupported canvas type")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrRateLimited is returned when the external system rate limits the request.
	ErrRateLimited = errors.New("rate limited")
)

// ProviderError wraps an error with provider context.
type ProviderError struct {
	Provider  string
	Operation string
	Err       error
}

func (e *ProviderError) Error() string {
	return e.Provider + ": " + e.Operation + ": " + e.Err.Error()
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// WrapError creates a ProviderError with context.
func WrapError(provider, operation string, err error) error {
	if err == nil {
		return nil
	}
	return &ProviderError{
		Provider:  provider,
		Operation: operation,
		Err:       err,
	}
}
