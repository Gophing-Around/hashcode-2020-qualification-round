package main

import (
	"fmt"
	"strings"
)

func getConfig(line string) (books int, libraries int, days int) {
	confParts := strings.Split(line, " ")
	books = toint(confParts[0])
	libraries = toint(confParts[1])
	days = toint(confParts[2])
	return
}

func buildBooks(scoreLine string, nBooks int) []Book {
	scores := strings.Split(scoreLine, " ")
	books := make([]Book, nBooks)

	for i := 0; i < nBooks; i++ {
		books[i] = Book{
			id:    i,
			score: toint(scores[i]),
		}
	}
	return books
}

func buildLibraries(lines []string, nLibraries int, availableBooks []Book) []*Library {
	libraries := make([]*Library, nLibraries)
	for i := 0; i < nLibraries*2; i += 2 {
		line1 := lines[i]
		line2 := lines[i+1]

		library := buildLibrary(line1, line2, availableBooks)
		library.id = fmt.Sprintf("%d", i/2)

		libraries[i/2] = library
	}
	return libraries
}

func buildLibrary(line1, line2 string, availableBooks []Book) *Library {
	libraryConfig := strings.Split(line1, " ")

	nBooks := toint(libraryConfig[0])
	library := Library{
		nBooks:        nBooks,
		signup:        toint(libraryConfig[1]),
		bookShippable: toint(libraryConfig[2]),
	}

	books := make([]Book, nBooks)
	for i, stringBookid := range strings.Split(line2, " ") {
		intBookID := toint(stringBookid)
		books[i] = availableBooks[intBookID]
	}
	library.books = books
	return &library
}
