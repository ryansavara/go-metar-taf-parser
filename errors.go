package metartafparser

import (
	"errors"
	"fmt"
)

// ParseError is the base error type for all parsing errors in this package.
type ParseError struct {
	message string
}

// NewParseError creates a new ParseError with the given message.
func NewParseError(message string) *ParseError {
	return &ParseError{message: message}
}

func (e *ParseError) Error() string { return e.message }

// InvalidWeatherStatementError indicates the input could not be parsed as a valid weather report.
type InvalidWeatherStatementError struct {
	ParseError
}

// NewInvalidWeatherStatementError returns a new InvalidWeatherStatementError, optionally with the invalid input as cause.
func NewInvalidWeatherStatementError(cause any) *InvalidWeatherStatementError {
	msg := "Invalid weather string"
	if s, ok := cause.(string); ok {
		msg = "Invalid weather string: " + s
	}
	return &InvalidWeatherStatementError{ParseError: ParseError{message: msg}}
}

func (e *InvalidWeatherStatementError) Unwrap() error { return &e.ParseError }

// PartialWeatherStatementError indicates the input appears to be a partial or incomplete TAF.
type PartialWeatherStatementError struct {
	InvalidWeatherStatementError

	Part  int
	Total int
}

// NewPartialWeatherStatementError returns a new PartialWeatherStatementError with details about the partial message.
func NewPartialWeatherStatementError(partialMessage string, part, total int) *PartialWeatherStatementError {
	return &PartialWeatherStatementError{
		InvalidWeatherStatementError: InvalidWeatherStatementError{
			ParseError: ParseError{
				message: fmt.Sprintf("Input is partial TAF (%s)", partialMessage),
			},
		},
		Part:  part,
		Total: total,
	}
}

func (e *PartialWeatherStatementError) Unwrap() error { return &e.InvalidWeatherStatementError }

var errCommandNotHandled = errors.New("command not handled")

type commandExecutionError struct {
	ParseError
}

func (e *commandExecutionError) Unwrap() error { return &e.ParseError }

// UnexpectedParseError indicates an unexpected error occurred during parsing.
type UnexpectedParseError struct {
	ParseError
}

// NewUnexpectedParseError creates a new UnexpectedParseError with the given message.
func NewUnexpectedParseError(message string) *UnexpectedParseError {
	return &UnexpectedParseError{ParseError: ParseError{message: message}}
}

func (e *UnexpectedParseError) Unwrap() error { return &e.ParseError }

// TimestampOutOfBoundsError indicates a parsed timestamp fell outside the expected range.
type TimestampOutOfBoundsError struct {
	ParseError
}

// NewTimestampOutOfBoundsError creates a new TimestampOutOfBoundsError with the given message.
func NewTimestampOutOfBoundsError(message string) *TimestampOutOfBoundsError {
	return &TimestampOutOfBoundsError{ParseError: ParseError{message: message}}
}

func (e *TimestampOutOfBoundsError) Unwrap() error { return &e.ParseError }
