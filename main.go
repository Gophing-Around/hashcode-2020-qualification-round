package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

// Library .
type Library struct {
	id            string
	nBooks        int
	signup        int
	bookShippable int
	books         []Book

	firstDayAvailable int

	sentBookIDs []string
	sentBooks   []Book

	libraryScore int
}

// Book .
type Book struct {
	id    int
	score int
}

func main() {
	files := []string{
		"a", // base
		"b", // 100k books | 100 libraries | 1000 d
		"c", // 100k books | 10k libraries | 100k d
		"d", // 78600 books | 30k libraries | 30001 d
		"e", // 100k books | 1k libraries | 200d
		"f", // 100k books | 1k libraries | 700d
	}

	for _, fileName := range files {
		fmt.Printf("****************** INPUT: %s\n", fileName)

		inputSet := readFile(fmt.Sprintf("./inputFiles/%s.in", fileName))

		configLines := strings.Split(inputSet, "\n")
		nBooks, nLibraries, nDays := getConfig(configLines[0])

		books := buildBooks(configLines[1], nBooks)
		libraries := buildLibraries(configLines[2:], nLibraries, books)

		// Sorting:
		//  - n giorni signup
		//  - n libri unici in libreria,
		//  - libri inviabili al giorno
		//  - score dei libri

		sortedLibraries := sortLibraries(libraries)

		outLibraries := algorithm(nDays, sortedLibraries, books)

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

func sortLibraries(libraries []*Library) []*Library {
	for _, lib := range libraries {
		bookShippable := lib.bookShippable
		uniqueBooks := lib.nBooks
		signup := lib.signup

		a := 1
		b := 1
		c := 1

		lib.libraryScore = ((bookShippable * a) * (uniqueBooks * b)) - (signup * c)
	}

	sort.Slice(libraries, func(i, j int) bool {
		libA := libraries[i]
		libB := libraries[j]
		return libA.libraryScore > libB.libraryScore
	})
	return libraries
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

func findLibrariesScanned(libraries []*Library) []*Library {
	newLibraries := make([]*Library, 0)
	for _, lib := range libraries {
		if len(lib.sentBooks) > 0 {
			newLibraries = append(newLibraries, lib)
		}
	}
	return newLibraries
}
