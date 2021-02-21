package main

import (
	"io/ioutil"
	"strconv"
)

func readFile(source string) string {
	in, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	return string(in)
}

func toint(num string) int {
	res, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	return res
}

func calcLibBookScore(books []*Book) (score int) {
	for _, book := range books {
		score += book.score
	}
	return
}

func removeElement(s []*Library, i int) []*Library {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
