package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

func parseArguments() string {
	csv := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	return *csv
}

func readFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	problems, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return problems, nil
}

func evaluateProblems(problems [][]string) int {
	var correctAnswers int
	for i, problem := range problems {
		var answer string

		fmt.Printf("Problem #%d: %s = ", i+1, problem[0])
		fmt.Scan(&answer)

		if answer == problem[1] {
			correctAnswers++
		}
	}
	return correctAnswers
}

func main() {
	csvFile := parseArguments()

	problems, err := readFile(csvFile)
	if err != nil {
		log.Fatal("ERROR ", err)
	}

	correctAnswers := evaluateProblems(problems)
	fmt.Printf("Your score is %d out of %d.\n", correctAnswers, len(problems))
}
