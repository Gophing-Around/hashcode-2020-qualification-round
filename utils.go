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
