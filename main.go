package main

import (
	"fmt"
	"net/http"
)

// Simple struct to hold any error that might occur processing a request
// and the associated http status code to return
type httpError struct {
	err        error
	statusCode int
}

// Custom Handler that can return errors
type appHandler func(w http.ResponseWriter, r *http.Request) *httpError

// Our custom Handler implements the http.Handler interface and wraps
// it with error handling capabilities
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.err.Error(), err.statusCode)
	}
}

func main() {
	http.Handle("/v1/users", appHandler(usersHandler))
	http.Handle("/v1/habits", appHandler(habitsHandler))
	http.Handle("/v1/habits/", appHandler(habitByIdHandler))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Started HabitTracker server on 8080.")
}
