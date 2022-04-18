package main

import (
	"bufio"
	"os"
	"strings"
)

type identifier interface {
	getID() int
}

// Takes:
// - A file name (will panic if does not exist in the current dir)
// - A transform function (from tokens, ie. an array/slice of string,
// to T, ie. something that implements the identifier interface above)
// and returns a map of T's keyed by that ID
func loadFromFile[T identifier](fileName string, transform func([]string) T) map[int]T {
	f, err := os.Open("./" + fileName)
	if err != nil {
		panic("File " + fileName + " not found")
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	res := make(map[int]T)
	for s.Scan() {
		tokens := strings.Split(s.Text(), ",")
		t := transform(tokens)
		res[t.getID()] = t
	}
	return res
}
