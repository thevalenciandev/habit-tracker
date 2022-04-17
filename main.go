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

// Custom Handler that can return errors
type appHandler func(w http.ResponseWriter, r *http.Request) error

// Our custom Handler implements the http.Handler interface and wraps
// it with error handling capabilities
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError) //TODO should we give more details of the error?
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

func habitByIdHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		matches := habitIdRegexp.FindStringSubmatch(r.URL.Path)
		if matches != nil || len(matches) != 2 {
			return errors.New("Path " + r.URL.Path + " not found")
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		h, ok := habits[id]
		if !ok {
			return errors.New("Habit " + strconv.Itoa(id) + " not found")
		} else {
			return encodeAsJson(h, w)
		}
	default:
		return errors.New("Method " + r.Method + " not supported")
	}
}

func encodeAsJson(toEncode any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(toEncode)
}

func allHabitsHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return encodeAsJson(maps.Values(habits), w)
	case "POST":
		var h Habit
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			return err
		}
		habitCnt++
		habits[habitCnt] = h
		return encodeAsJson(h, w)
	default:
		return errors.New("Method " + r.Method + " not supported")
	}
}
