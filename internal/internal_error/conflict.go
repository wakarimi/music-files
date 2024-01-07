package internal_error

import "fmt"

type Conflict struct {
	Message string
}

func (e Conflict) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
