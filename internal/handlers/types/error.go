package types

// ErrorResponse godoc
// @Description Standard error response
// @Property ErrorResponse (string, required) A human-readable description of the error
type ErrorResponse struct {
	Error string `json:"error" binding:"required"`
}
