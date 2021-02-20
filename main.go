package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files := []string{
		"aple",
		"b.in",
		"c.in",
		"d.in",
		"e.in",
		"f.in",
	}

	for _, fileName := range files {
		fmt.Printf("--------------------------------------------------------")
		fmt.Printf("****************** INPUT: %s\n", fileName)

		inputSet := readFile(fmt.Sprintf("./inputFiles/%s", fileName))

		// Algorithm here
		// output := algorithm()

		result := fmt.Sprintf("%d\n", len(output))
		for _, order := range output {
			// result += fmt.Sprintf("%d %s\n", len(order.pizzas), strings.Join(order.pizzas, " "))
		}
		ioutil.WriteFile(fmt.Sprintf("./result/%s", fileName), []byte(result), 0644)
	}
}

func algorithm() {}
