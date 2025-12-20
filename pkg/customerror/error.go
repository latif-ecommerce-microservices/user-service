package customerror

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

type Error struct {
	code       Code
	message    string
	metadata   map[string]interface{}
	locator    *Locator
	baseErr    error
	stackTrace pkgerrors.StackTrace
}

func (ce Error) Error() string {
	return ce.message
}

func (ce Error) Code() string {
	return string(ce.code)
}

func (ce Error) Metadata() map[string]interface{} {
	return ce.metadata
}

func (ce Error) WithMetadata(metadata map[string]interface{}) Error {
	ce.metadata = metadata
	return ce
}

func (ce Error) Sprintf(a ...any) Error {
	newCe := ce
	newCe.message = fmt.Sprintf(ce.message, a...)
	return newCe
}

func (ce Error) WithLocator(location *Locator) Error {
	ce.locator = location
	return ce
}

func (ce Error) WithCause(err error) Error {
	ce.baseErr = err
	return ce
}

func (ce Error) BaseError() error {
	return ce.baseErr
}

func (ce Error) Locator() *Locator {
	return ce.locator
}

func (ce Error) WithStackTrace() Error {
	if ce.baseErr == nil {
		return ce
	}

	type stackTracer interface {
		StackTrace() pkgerrors.StackTrace
	}

	if st, ok := ce.baseErr.(stackTracer); ok {
		ce.stackTrace = st.StackTrace()
	} else {
		ce.baseErr = pkgerrors.WithStack(ce.baseErr)
		if st, ok := ce.baseErr.(stackTracer); ok {
			ce.stackTrace = st.StackTrace()
		}
	}

	return ce
}

func (ce Error) StackTrace() pkgerrors.StackTrace {
	return ce.stackTrace
}

func NewErrorWithCode(code Code, message string) Error {
	return Error{code, message, nil, nil, nil, nil}
}

func NewErrorWithCodeAndLocator(code Code, message string, locator *Locator) Error {
	return Error{code, message, nil, locator, nil, nil}
}

func NewClientError(code Code, message string) Error {
	return NewErrorWithCode(code, message)
}

func NewForbiddenError(code Code, message string) Error {
	return NewErrorWithCode(code, message)
}

func NewUnauthorizedError(code Code, message string) Error {
	return NewErrorWithCode(code, message)
}

func NewTooManyRequestError(code Code, message string) Error {
	return NewErrorWithCode(code, message)
}

func NewNotFoundError(code Code, message string) Error {
	return NewErrorWithCode(code, message)
}

func NewInternalServerError(code Code, message string, locator ...Locator) Error {
	if len(locator) == 0 {
		return NewErrorWithCode(code, message)
	}

	return NewErrorWithCodeAndLocator(code, message, &locator[0])
}
