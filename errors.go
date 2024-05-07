package slogerror

import (
	"fmt"
	"log/slog"
)

type slogError struct {
	err     string
	wrapped error
}

func (se slogError) Error() string {
	return se.err
}

func (se slogError) Unwrap() error {
	return se.wrapped
}

// Error returns an error with context provided by
// the passed logger.
func Error(l *slog.Logger, ob any) slogError {
	return slogError{err: contextString(l) + " " + fmt.Sprintf("%v", ob)}
}

// Wrap wraps the provided error, includes the passed object
// in the message and includes context from the passed logger.
func Wrap(l *slog.Logger, ob any, err error) slogError {
	return slogError{
		err:     contextString(l) + " " + fmt.Sprintf("%v: %s", ob, err.Error()),
		wrapped: err,
	}
}

func contextString(l *slog.Logger) string {
	h, isType := l.Handler().(*Handler)
	if !isType {
		return "[logger not using slogerror]"
	}
	buf := make([]byte, 0, 1024)
	for _, attr := range h.attrs {
		buf = h.contextString(buf, "", attr)
	}
	return string(buf)
}
