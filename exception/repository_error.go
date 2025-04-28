package exception

import "fmt"

type RepositoryError struct {
	Message string
	FunctionName string
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("[Repository Error][%s] %s", e.FunctionName, e.Message)
}

func NewRepositoryError(message string, functionName string) *RepositoryError {
	return &RepositoryError{
		Message: message,
		FunctionName: functionName,
	}
}
