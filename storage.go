package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func loadUsers() map[int]User {
	f, err := os.Open("./users.csv")
	if err != nil {
		panic("File users.csv not found")
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	res := make(map[int]User)
	for s.Scan() {
		tokens := strings.Split(s.Text(), ",")
		id, err := strconv.Atoi(strings.TrimSpace(tokens[0]))
		if err != nil {
			panic("Error loading users. Invalid ID.")
		}
		name := strings.TrimSpace(tokens[1])
		res[id] = User{id, name}
	}
	return res
}
