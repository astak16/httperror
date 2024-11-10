package httperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Err    error
	Status int
}

func Wrap(err error, status int) error {
	if err == nil {
		return nil
	}
	return Error{err, status}
}

func Errorf(status int, format string, args ...interface{}) error {
	return Wrap(fmt.Errorf(format, args...), status)
}

func (e Error) Error() string {
	return e.Err.Error()
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, status int)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
}

func NewWithHandler(next Handler, eh ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next.ServeHTTP(w, r)
		if err == nil {
			return
		}

		herr := Error{}
		if errors.As(err, &herr) {
			eh(w, r, herr, herr.Status)
		} else {
			eh(w, r, err, http.StatusInternalServerError)
		}
	})
}

func New(next Handler) http.Handler {
	return NewWithHandler(next, DefaultErrorHandler)
}

func NewF(next HandlerFunc) http.Handler {
	return New(next)
}

type errorResponse struct {
	Message string `json:"message"`
}
