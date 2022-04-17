package main

import (
	"encoding/json"
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

var (
	habitIdRegexp = regexp.MustCompile(`/v1/habits/([0-9]+)$`)
	habits        = make(map[int]Habit)
	habitCnt      = 1
)

func main() {
	habits[habitCnt] = Habit{Id: "1", Name: "Play Guitar 🎸", Days: []string{"Monday", "Sunday"}}

	http.HandleFunc("/v1/habits", allHabitsHandler)
	http.HandleFunc("/v1/habits/", habitByIdHandler)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Started HabitTracker server on 8080.")
}

func habitByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		matches := habitIdRegexp.FindStringSubmatch(r.URL.Path)
		if matches != nil && len(matches) == 2 {
			id, err := strconv.Atoi(matches[1])
			if err == nil {
				fmt.Println("Id requested", id)
				h, ok := habits[id]
				if ok {
					encodeAsJson(h, w)
					return
				} else {
					http.Error(w, "Not found", http.StatusNotFound)
					return
				}
			}
		}
	}
	http.Error(w, "Unexpected error", http.StatusInternalServerError) //TODO should we give more details of the error?
}

func encodeAsJson(toEncode any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(toEncode)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
	}
	return err
}

func allHabitsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		encodeAsJson(maps.Values(habits), w)
	case "POST":
		var h Habit
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			http.Error(w, "Unexpected error decoding habit", http.StatusNoContent)
			return
		}
		habitCnt++
		habits[habitCnt] = h
		encodeAsJson(h, w)
	default:
		http.Error(w, "Method not supported ", http.StatusMethodNotAllowed)
	}
}
