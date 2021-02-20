package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	files := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
	}

	for _, fileName := range files {
		fmt.Printf("--------------------------------------------------------")
		fmt.Printf("****************** INPUT: %s\n", fileName)

		inputSet := readFile(fmt.Sprintf("./inputFiles/%s.in", fileName))

		configLines := strings.Split(inputSet, "\n")
		nBooks, nLibraries, nDays := getConfig(configLines[0])

		books := buildBooks(configLines[1], nBooks)
		libraries := buildLibraries(configLines[2:], nLibraries, books)

		outLibraries := algorithm(nDays, libraries, books)

		scannedLibraries := findLibrariesScanned(outLibraries)

		result := fmt.Sprintf("%d\n", len(scannedLibraries))
		for _, lib := range scannedLibraries {
			result += fmt.Sprintf("%s %d\n", lib.id, len((lib.sentBooks)))
			result += fmt.Sprintf("%s\n", strings.Join(lib.sentBookIDs, " "))
		}

		result = strings.TrimSpace(result)
		ioutil.WriteFile(fmt.Sprintf("./result/%s.out", fileName), []byte(result), 0644)
	}
}

func findLibrariesScanned(libraries []*Library) []*Library {
	newLibraries := make([]*Library, 0)
	for _, lib := range libraries {
		if len(lib.sentBooks) > 0 {
			newLibraries = append(newLibraries, lib)
		}
	}
	return newLibraries
}

func algorithm(nDays int, libraries []*Library, books []Book) []*Library {
	sentbooks := make(map[int]bool)

	startingDay := 0
	for _, library := range libraries {
		library.firstDayAvailable = startingDay + library.signup
		startingDay += library.signup
	}

	for day := 0; day < nDays; day++ {
		for _, library := range libraries {
			if library.firstDayAvailable < day {
				continue
			}

			shippablePerLibrary := library.bookShippable

			for _, book := range library.books {
				if shippablePerLibrary == 0 {
					break
				}

				if sent, ok := sentbooks[book.id]; !sent || !ok {
					library.sentBooks = append(library.sentBooks, book)
					library.sentBookIDs = append(library.sentBookIDs, fmt.Sprintf("%d", book.id))
					sentbooks[book.id] = true
					shippablePerLibrary--
				}
			}
		}
	}

	return libraries
}

type Library struct {
	id            string
	nBooks        int
	signup        int
	bookShippable int
	books         []Book

	firstDayAvailable int

	sentBookIDs []string
	sentBooks   []Book
}

type Book struct {
	id    int
	score int
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

func getConfig(line string) (books int, libraries int, days int) {
	confParts := strings.Split(line, " ")
	books = toint(confParts[0])
	libraries = toint(confParts[1])
	days = toint(confParts[2])
	return
}
