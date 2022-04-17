package main

import (
	"encoding/json"
	"errors"
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

func allHabitsHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		err := encodeAsJson(maps.Values(habits), w, http.StatusOK)
		return err
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
