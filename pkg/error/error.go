package error

// ErrorResponse is a type of response for when the error comes from the user.
type ErrorResponse struct {
	// ErrorMessage define the ErrorMessage message returned by the api
	ErrorMessage string `example:"this ingredient doesn't exist"`
}
