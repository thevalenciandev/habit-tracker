package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/exp/maps"
)

// Public as the encoding/json package needs access to it to do its job
type Habit struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Days []string `json:"days"`
}

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

var (
	habitIdRegexp = regexp.MustCompile(`/v1/habits/([0-9]+)$`)
	habits        = make(map[int]Habit)
	habitCnt      = 1
)

func main() {
	habits[habitCnt] = Habit{Id: "1", Name: "Play Guitar 🎸", Days: []string{"Monday", "Sunday"}}

	http.Handle("/v1/habits", appHandler(allHabitsHandler))
	http.Handle("/v1/habits/", appHandler(habitByIdHandler))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Started HabitTracker server on 8080.")
}

func habitByIdHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		matches := habitIdRegexp.FindStringSubmatch(r.URL.Path)
		if matches == nil || len(matches) != 2 {
			return &httpError{errors.New("Path " + r.URL.Path + " not found"), http.StatusNotFound}
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			return &httpError{err, http.StatusNotFound}
		}
		h, ok := habits[id]
		if !ok {
			return &httpError{errors.New("Habit " + strconv.Itoa(id) + " not found"), http.StatusNotFound}
		}
		return encodeAsJson(h, w, http.StatusOK)

	default:
		return &httpError{errors.New("Method " + r.Method + " not supported"), http.StatusMethodNotAllowed}
	}
}

func encodeAsJson(toEncode any, w http.ResponseWriter, statusCode int) *httpError {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toEncode); err != nil {
		return &httpError{err, http.StatusInternalServerError}
	}
	return nil
}

func allHabitsHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		return encodeAsJson(maps.Values(habits), w, http.StatusOK)
	case "POST":
		var h Habit
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			return &httpError{err, http.StatusInternalServerError}
		}
		habitCnt++
		habits[habitCnt] = h
		return encodeAsJson(h, w, http.StatusCreated)
	default:
		return &httpError{errors.New("Method " + r.Method + " not supported"), http.StatusMethodNotAllowed}
	}
}
