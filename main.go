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

	sendAttempts int
	sentBookIDs  []string
	sentBooks    []Book

	libraryScore int
}

// Book .
type Book struct {
	id    int
	score int
}

func main() {
	files := []string{
		// "a", // base
<<<<<<< HEAD
		// "b", // 100k books | 100 libraries | 1000 days
		"c", // 100k books | 10k libraries | 100k days
		// "d", // 78600 books | 30k libraries | 30001 days
		// "e", // 100k books | 1k libraries | 200 days
		// "f", // 100k books | 1k libraries | 700 days
=======
		"b", // 100k books | 100 libraries | 1000 days
		// "c", // 100k books | 10k libraries | 100k days
		// "d", // 78600 books | 30k libraries | 30001 days
		"e", // 100k books | 1k libraries | 200 days
		"f", // 100k books | 1k libraries | 700 days
>>>>>>> 9e8d0589aac2caa0556c69e43062c943660f19d0
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

		fmt.Printf("Scanned libraries: %d - Total libraries: %d\n", len(scannedLibraries), len(libraries))

		result := fmt.Sprintf("%d\n", len(scannedLibraries))
		mean := 0.0
		for _, lib := range scannedLibraries {
			// fmt.Printf("Sent books for library %s: %d\n", lib.id, len(lib.sentBooks))
			total := lib.nBooks
			sent := len(lib.sentBooks)
			percent := sent / total
			mean = (mean + float64(percent)) / float64(len(scannedLibraries))

			result += fmt.Sprintf("%s %d\n", lib.id, len((lib.sentBooks)))
			result += fmt.Sprintf("%s\n", strings.Join(lib.sentBookIDs, " "))
		}

		fmt.Printf("Sent book mean: %.000f\n", mean)

		result = strings.TrimSpace(result)
		ioutil.WriteFile(fmt.Sprintf("./result/%s.out", fileName), []byte(result), 0644)
	}
}

func sortLibraries(libraries []*Library) []*Library {
	for _, lib := range libraries {
		bookShippable := lib.bookShippable
		nbooks := lib.nBooks
		libraryBooksScore := calcLibBookScore(lib.books) / nbooks
		signupDays := lib.signup

		bookShippableCoef := 1
		libraryBooksScoreCoef := 100
		signupDaysCoef := 1000

		lib.libraryScore = ((bookShippable * bookShippableCoef) *
			(libraryBooksScore * libraryBooksScoreCoef)) /
			(signupDays * signupDaysCoef)
	}

	sort.Slice(libraries, func(i, j int) bool {
		libA := libraries[i]
		libB := libraries[j]
		return libA.libraryScore > libB.libraryScore
	})
	return libraries
}

func updateLibraryScores(libraries []*Library, sentbooks map[int]bool) []*Library {
	uniqueBooksAvailableCoef := 1

	for _, lib := range libraries {
		uniqueBooksAvailable := 0

		for _, book := range lib.books {
			if sent, ok := sentbooks[book.id]; !sent || !ok {
				uniqueBooksAvailable++
			}
		}

		lib.libraryScore *= uniqueBooksAvailable * uniqueBooksAvailableCoef
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

		libraries = updateLibraryScores(libraries, sentbooks)

		for _, library := range libraries {
			if library.firstDayAvailable < day {
				continue
			}

			if library.sendAttempts == library.nBooks {
				continue
			}

			shippablePerLibrary := library.bookShippable
			for _, book := range library.books {
				if shippablePerLibrary == 0 {
					break
				}

				library.sendAttempts++
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
