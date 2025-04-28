package exception

import "fmt"

type ServiceError struct {
	Message string
	FunctionName string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("[Service Error][%s] %s", e.FunctionName, e.Message )
}

func NewServiceError(message string, functionName string) *ServiceError {
	return &ServiceError{
		Message: message,
		FunctionName: functionName,
	}
}
