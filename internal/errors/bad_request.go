package errors

import "fmt"

type BadRequest struct {
	Message string
}

func (e BadRequest) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
