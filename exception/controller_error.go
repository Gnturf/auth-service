package exception

import "fmt"

type ControllerError struct {
	Message string
	FunctionName string
}

func (e *ControllerError) Error() string {
	return fmt.Sprintf("[Service Error][%s] %s", e.FunctionName, e.Message )
}

func NewControllerError(message string, functionName string) *ControllerError {
	return &ControllerError{
		Message: message,
		FunctionName: functionName,
	}
}
