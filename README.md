# errors
--
    import "github.com/qri-io/errors"

Package errors defines a rich error type with categorization codes, friendly
messages, support for a lower-level causer interface, and presentational error
utility functions

this error package works best when used sparingly at critical points in an
application stack, only adding context & detail when handling the error is in a
position to provide additional error details for the end user

any lower level error would be better off using github.com/pkg/errors, which
will interoperate nicely with this package

## Usage

#### func  Cause

```go
func Cause(err error) error
```
Cause proxies the pkg/errors Cause function

#### func  CodeHTTPStatus

```go
func CodeHTTPStatus(c Code) int
```
CodeHTTPStatus converts a Code to an http status code, defaulting to 500

#### func  CodeString

```go
func CodeString(c Code) string
```
CodeString returns a string representation of a code, defaulting to "error"

#### func  RegisterCode

```go
func RegisterCode(c Code, httpStatus int, typeStr string) error
```
RegisterCode adds a code to error's internal code pool for extending Error with
custom http and string values for codes

#### type Code

```go
type Code int
```

Code assigns numeric values to different categories of error Users are
encouraged to define their own application-specific errors

```go
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
```

#### type Error

```go
type Error struct {
}
```

Error decorates an error with additional fields for user feedback It couples
developer-focused errors with a code for classifying the error, and an optional
user-friendly error message. values that caused the error to occur should be
given to the error as data params

Errors should always be created with New or one of it's variants

#### func  New

```go
func New(c Code, message string, data ...interface{}) *Error
```
New creates an Error from an error and string

#### func  NewFriendly

```go
func NewFriendly(c Code, message, friendly string, data ...interface{}) *Error
```
NewFriendly creates an error with a user-friendly message

#### func  NewFriendlyFix

```go
func NewFriendlyFix(c Code, message, friendly, fix string, data ...interface{}) *Error
```
NewFriendlyFix creates an error with a message and a fix

#### func  Wrap

```go
func Wrap(c Code, err error, message string, data ...interface{}) *Error
```
Wrap returns an error annotating err with a stack trace at the point Wrap is
called, and the supplied message. If err is nil, Wrap returns nil.

#### func  WrapFriendly

```go
func WrapFriendly(c Code, err error, message, friendly string, data ...interface{}) *Error
```
WrapFriendly calls wrap and adds a friendly, user-facing message describing the
problem

#### func  WrapFriendlyFix

```go
func WrapFriendlyFix(c Code, err error, message, friendly, fix string, data ...interface{}) *Error
```
WrapFriendlyFix calls wrap and adds a friendly, a user-facing message describing
the problem

#### func (Error) Cause

```go
func (e Error) Cause() error
```
Cause implements the causer interface from the errors standard package

#### func (Error) Code

```go
func (e Error) Code() Code
```
Code gives the type of error

#### func (Error) Error

```go
func (e Error) Error() string
```
Error satisfies the error interface, printing just top-level error

#### func (Error) Fix

```go
func (e Error) Fix() string
```
Fix returns the internal message on how to fix the error

#### func (Error) Friendly

```go
func (e Error) Friendly() string
```
Friendly returns the friendly message along
