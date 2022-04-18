package main

import (
	"net/http"

	"golang.org/x/exp/maps"
)

type User struct {
	name string
}

var users = make(map[int]User)

func usersHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		return encodeAsJson(maps.Values(users), w, http.StatusOK)
	default:
		return methodNotAllowedHttpError(r.Method)
	}
}
