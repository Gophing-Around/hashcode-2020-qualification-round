package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Library .
type Library struct {
	id            string
	nBooks        int
	signup        int
	bookShippable int
	books         []*Book

	firstDayAvailable int

	sendAttempts int
	sentBookIDs  []string
	sentBooks    []*Book

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
		"b", // 100k books | 100 libraries | 1000 days
		"e", // 100k books | 1k libraries | 200 days
		"f", // 100k books | 1k libraries | 700 days
		"c", // 100k books | 10k libraries | 100k days
		"d", // 78600 books | 30k libraries | 30001 days
	}

	for _, fileName := range files {
		fmt.Printf("****************** INPUT: %s\n", fileName)

		inputSet := readFile(fmt.Sprintf("./inputFiles/%s.in", fileName))

		configLines := strings.Split(inputSet, "\n")
		nBooks, nLibraries, nDays := getConfig(configLines[0])

		books := buildBooks(configLines[1], nBooks)
		libraries := buildLibraries(configLines[2:], nLibraries, books)

		outLibraries := algorithm(nDays, libraries, books)
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

func findBestLibrary(libraries []*Library, remainingDays int, sentbooks map[int]bool) (int, *Library) {
	booksScoreCoef := 1000
	signupDaysCoef := 10
	wastedTimeCoef := 10

	maxScore, maxScoreIndex := -1, -1
	for index, library := range libraries {
		if library.signup > remainingDays {
			continue
		}

		nBooksToTake := ( remainingDays - library.signup ) * library.bookShippable
		count := 0
		score := 0
		for _, book := range library.books {
			if count >= nBooksToTake {
				break
			}

			if sent, ok := sentbooks[book.id]; !sent || !ok {
				score += book.score
				count++
			}
		}

		wastedTimePenalty := nBooksToTake - count
		if wastedTimePenalty <= 0 {
			wastedTimePenalty = 1
		} else {
			wastedTimePenalty *= wastedTimeCoef
		}

		score = (score * booksScoreCoef)	/	((library.signup * signupDaysCoef)) // * wastedTimePenalty) // wasted time seems to not work
		if score > maxScore {
			maxScore = score
			maxScoreIndex = index
		}
	}

	if maxScore < 0 || maxScoreIndex < 0 {
		return -1, nil
	}

	return maxScoreIndex, libraries[maxScoreIndex]
}

func algorithm(nDays int, origLibraries []*Library, books []*Book) []*Library {
	sentbooks := make(map[int]bool)
	signedUpLibraries := make([]*Library, 0)
	libraries := origLibraries

	var currentSignignLibrary *Library = nil
	indexToRemove := 0
	lastSigningStartingDay := 0
	for day := 0; day < nDays; day++ {
		for _, library := range signedUpLibraries {
			if library.sendAttempts >= library.nBooks {
				continue
			}

			library.sendAttempts = 0
			shippablePerLibrary := (nDays-day)*library.bookShippable
			for _, book := range library.books {
				if shippablePerLibrary <= 0 {
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

		if currentSignignLibrary == nil {
			if len(libraries) > 0 {
				indexToRemove, currentSignignLibrary = findBestLibrary(libraries, nDays-day, sentbooks)
				if indexToRemove > 0 && currentSignignLibrary != nil {
					libraries = removeElement(libraries, indexToRemove)
					lastSigningStartingDay = day
				}
			}
		} else if (day - lastSigningStartingDay >= currentSignignLibrary.signup-1) {
			signedUpLibraries = append(signedUpLibraries, currentSignignLibrary)
			currentSignignLibrary = nil
		}
	}

	return origLibraries
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
