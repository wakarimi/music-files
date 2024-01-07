package internal_error

import "fmt"

type NotFound struct {
	EntityName string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.EntityName)
}
