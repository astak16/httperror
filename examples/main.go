package main

import (
	"fmt"
	"github.com/astak16/httperror"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/s", httperror.NewF(func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(w, "this is OK")
		return nil
	}))

	mux.Handle("/f", httperror.NewF(func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("this will be a 500")
	}))

	mux.Handle("/e", httperror.NewF(func(w http.ResponseWriter, r *http.Request) error {
		return httperror.Wrap(fmt.Errorf("wrap another err into a bad request"), http.StatusBadRequest)
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
