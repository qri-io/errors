// Package errors defines a rich error type with categorization codes, friendly
// messages, support for a lower-level causer interface, and presentational
// error utility functions
//
// this error package works best when used sparingly at critical points in
// an application stack, only adding context & detail when handling the error
// is in a position to provide additional error details for the end user
//
// any lower level error would be better off using github.com/pkg/errors, which
// will interoperate nicely with this package
package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Code assigns numeric values to different categories of error
// Users are encouraged to define their own application-specific errors
type Code int

const (
	// CodeUnknown should never be used, indicates unspecified default
	CodeUnknown Code = iota
	// CodeGeneric indices a nonspecific error. use this if you're not sure what
	// error code to use
	CodeGeneric
	// CodeInvalidSyntax indicates provided values cannot be deserialized
	CodeInvalidSyntax
	// CodeInvalidArgs indicates the provided arguments are invalid
	CodeInvalidArgs
	// CodeUnauthorized there is a problem with the client’s credentials
	CodeUnauthorized
	// CodeForbidden indicates access isn't permitted, regardless of
	// authorization state
	CodeForbidden
	// CodeNotFound indicates a provided value couldn't be retrieved
	CodeNotFound
	// CodeUnavailable indicates something that needs to be available cannot
	// be reached
	CodeUnavailable
)

type codeDetails struct {
	httpStatus int
	str        string
}

var codePool = map[Code]codeDetails{
	CodeUnknown:       {500, "error"},
	CodeGeneric:       {500, "error"},
	CodeInvalidSyntax: {400, "syntax"},
	CodeInvalidArgs:   {400, "arguments"},
	CodeUnauthorized:  {401, "auth"},
	CodeForbidden:     {403, "auth"},
	CodeNotFound:      {404, "missing"},
	CodeUnavailable:   {503, "unavailable"},
}

// RegisterCode adds a code to error's internal code pool for extending Error with
// custom http and string values for codes
func RegisterCode(c Code, httpStatus int, typeStr string) error {
	if _, ok := codePool[c]; ok {
		return New(CodeInvalidArgs, "already registered", c)
	}
	codePool[c] = codeDetails{httpStatus, typeStr}
	return nil
}

// CodeString returns a string representation of a code, defaulting to "error"
func CodeString(c Code) string {
	if s, ok := codePool[c]; ok {
		return s.str
	}
	return "error"
}

// CodeHTTPStatus converts a Code to an http status code, defaulting to 500
func CodeHTTPStatus(c Code) int {
	if s, ok := codePool[c]; ok {
		return s.httpStatus
	}
	return 500
}

// Error decorates an error with additional fields for user feedback
// It couples developer-focused errors with a code for classifying the error,
// and an optional user-friendly error message. values that caused the error
// to occur should be given to the error as data params
//
// Errors should always be created with New or one of it's variants
type Error struct {
	code     Code
	friendly string
	fix      string
	data     []interface{}
	cause    error
}

// Error satisfies the error interface, printing just top-level error
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", CodeString(e.code), e.cause.Error())
}

// Cause implements the causer interface from the errors standard package
func (e Error) Cause() error {
	return e.cause
}

// Code gives the type of error
func (e Error) Code() Code {
	return e.code
}

// Fix returns the internal message on how to fix the error
func (e Error) Fix() string {
	return e.fix
}

// Friendly returns the friendly message along
func (e Error) Friendly() string {
	if e.friendly == "" && e.fix == "" {
		return ""
	}

	str := fmt.Sprintf("%s: %s", CodeString(e.code), e.friendly)
	for i, d := range e.data {
		str += fmt.Sprintf(" %v", d)
		if i < len(e.data)-1 {
			str += ","
		} else {
			str += "."
		}
	}
	if e.fix != "" {
		str += fmt.Sprintf(" %s", e.fix)
	}
	return str
}

// New creates an Error from an error and string
func New(c Code, message string, data ...interface{}) *Error {
	err := errors.New(message)
	return &Error{code: c, data: data, cause: err}
}

// NewFriendly creates an error with a user-friendly message
func NewFriendly(c Code, message, friendly string, data ...interface{}) *Error {
	err := New(c, message, data...)
	err.friendly = friendly
	return err
}

// NewFriendlyFix creates an error with a message and a fix
func NewFriendlyFix(c Code, message, friendly, fix string, data ...interface{}) *Error {
	err := NewFriendly(c, message, friendly, data...)
	err.fix = fix
	return err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(c Code, err error, message string, data ...interface{}) *Error {
	return &Error{code: c, data: data, cause: errors.Wrap(err, message)}
}

// WrapFriendly calls wrap and adds a friendly, user-facing message describing the problem
func WrapFriendly(c Code, err error, message, friendly string, data ...interface{}) *Error {
	e := Wrap(c, err, message, data...)
	e.friendly = friendly
	return e
}

// WrapFriendlyFix calls wrap and adds a friendly, a user-facing message describing the problem
func WrapFriendlyFix(c Code, err error, message, friendly, fix string, data ...interface{}) *Error {
	e := WrapFriendly(c, err, message, friendly, data...)
	e.fix = fix
	return e
}

// Cause proxies the pkg/errors Cause function
func Cause(err error) error {
	return errors.Cause(err)
}
