package model

type Error interface {
	Error() string
	ErrorCode() int
}
type ApiError struct {
	Err  string
	Code int
}

func (e *ApiError) Error() string {
	return e.Err
}
func (e *ApiError) ErrorCode() int {
	return -1000
}
