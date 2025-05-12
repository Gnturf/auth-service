package exception

import "fmt"

type RepositoryError struct {
	Message string
	OriginalError error
	FunctionName string
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("[Repository Error][%s] %s", e.FunctionName, e.Message)
}

func (e *RepositoryError) Unwrap() error {
	return e.OriginalError
}

func NewRepositoryError(message string, originalError error, functionName string) *RepositoryError {
	return &RepositoryError{
		Message: message,
		FunctionName: functionName,
		OriginalError: originalError,
	}
}
