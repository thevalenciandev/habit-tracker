package main

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type User struct {
	ID   int
	Name string
}

// Implement the getID contract from storage
func (u User) getID() int {
	return u.ID
}

var users map[int]User = loadFromFile("users.csv", userTransform)

func userTransform(fileLineTokens []string) User {
	id, err := strconv.Atoi(strings.TrimSpace(fileLineTokens[0]))
	if err != nil {
		panic("Error loading user. Invalid ID: " + fileLineTokens[0])
	}
	name := strings.TrimSpace(fileLineTokens[1])
	return User{id, name}
}

func usersHandler(w http.ResponseWriter, r *http.Request) *httpError {
	switch r.Method {
	case "GET":
		return encodeAsJson(maps.Values(users), w, http.StatusOK)
	default:
		return methodNotAllowedHttpError(r.Method)
	}
}
