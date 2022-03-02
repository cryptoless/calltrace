package model

import "fmt"

type methodNotFoundError struct {
	Method string
}

func (e *methodNotFoundError) Error() string {
	return fmt.Sprintf("the method %s does not exist/is not available", e.Method)
}
func (e *methodNotFoundError) ErrorCode() int {
	return -1001
}

func MethodNotFoundError(method string) *methodNotFoundError {
	return &methodNotFoundError{
		Method: method,
	}
}

type parseError struct {
	Msg string
}

func (e *parseError) Error() string {
	return e.Msg
}
func (e *parseError) ErrorCode() int {
	return -1002
}
func ParseError(err string) *parseError {
	return &parseError{
		Msg: err,
	}
}

type invalidRequestError struct {
	Msg string
}

func (e *invalidRequestError) Error() string {
	return e.Msg
}
func (e *invalidRequestError) ErrorCode() int {
	return -1002
}
func InvalidRequestError(err string) *invalidRequestError {
	return &invalidRequestError{
		Msg: err,
	}
}

type invalidMessageError struct {
	Msg string
}

func (e *invalidMessageError) Error() string {
	return e.Msg
}
func (e *invalidMessageError) ErrorCode() int {
	return -1004
}
func InvalidMessageError(err string) *invalidMessageError {
	return &invalidMessageError{
		Msg: err,
	}
}

type invalidParamsError struct {
	Msg string
}

func (e *invalidParamsError) Error() string {
	return e.Msg
}
func (e *invalidParamsError) ErrorCode() int {
	return -1005
}
func InvalidParamsError(err string) *invalidParamsError {
	return &invalidParamsError{
		Msg: err,
	}
}

type internalError struct {
	Msg string
}

func (e *internalError) Error() string {
	return e.Msg
}
func (e *internalError) ErrorCode() int {
	return -1006
}
func InternalError(err string) *internalError {
	return &internalError{
		Msg: err,
	}
}

type subscriptionError struct {
	Method string
}

func (e *subscriptionError) Error() string {
	return fmt.Sprintf("no %q subscription in namespace", e.Method)
}
func (e *subscriptionError) ErrorCode() int {
	return -1006
}
func SubscriptionError(mothod string) *subscriptionError {
	return &subscriptionError{
		Method: mothod,
	}
}
