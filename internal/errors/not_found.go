package errors

import "fmt"

type NotFound struct {
	Resource string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}
