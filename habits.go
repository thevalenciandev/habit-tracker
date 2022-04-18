package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

// Public as the encoding/json package needs access to it to do its job
type Habit struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Days []string `json:"days"`
}

// Implement the getID contract from storage
func (h Habit) getID() int {
	return h.ID
}

var (
	habitIdRegexp               = regexp.MustCompile(`/v1/habits/([0-9]+)$`)
	habits        map[int]Habit = loadFromFile("habits.csv", habitsTransform)
	habitCnt                    = len(habits) + 1 // assumes IDs start from 1
)

// Assumes Habits are stored as ID,Name,Days
// where Days are a "|" separated list of weekdays
func habitsTransform(fileLineTokens []string) Habit {
	id, err := strconv.Atoi(strings.TrimSpace(fileLineTokens[0]))
	if err != nil {
		panic("Error loading habit. Invalid ID: " + fileLineTokens[0])
	}
	name := strings.TrimSpace(fileLineTokens[1])
	days := strings.Split(strings.TrimSpace(fileLineTokens[2]), "|")
	return Habit{id, name, days}
}

func habitsHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		err := encodeAsJson(maps.Values(habits), w, http.StatusOK)
		return err
	case "POST":
		var h Habit
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			return &httpError{err, http.StatusInternalServerError}
		}
		h.ID = habitCnt
		habits[habitCnt] = h
		habitCnt++
		return encodeAsJson(h, w, http.StatusCreated)
	default:
		return methodNotAllowedHttpError(r.Method)
	}
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
		return methodNotAllowedHttpError(r.Method)
	}
}
