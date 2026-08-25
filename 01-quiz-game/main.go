package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

func parseArguments() (csv string, limit int, shuffle bool) {
	flag.StringVar(&csv, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle the quiz order")
	flag.Parse()

	return csv, limit, shuffle
}

func parseCSV(filename string) ([]Problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf(`Failed to open the file: "%s"`, filename)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the provided CSV file.")
	}

	problems := make([]Problem, len(lines))
	for i, line := range lines {
		problems[i] = Problem{
			Question: line[0],
			Answer:   strings.TrimSpace(line[1]),
		}
	}

	return problems, nil
}

func readAnswer(number int, question string, c chan<- string) {
	fmt.Printf("Problem #%d: %s = ", number, question)
	var answer string
	fmt.Scan(&answer)
	c <- answer
}

func evaluateProblems(problems []Problem, timer *time.Timer) int {
	var correct int
	answerCh := make(chan string)
	for i, p := range problems {
		go readAnswer(i+1, p.Question, answerCh)

		select {
		case <-timer.C:
			fmt.Println("")
			return correct
		case answer := <-answerCh:
			if strings.EqualFold(answer, p.Answer) {
				correct++
			}
		}
	}

	return correct
}

func main() {
	csvFilename, timeLimit, shuffle := parseArguments()

	problems, err := parseCSV(csvFilename)
	if err != nil {
		log.Fatal("ERROR ", err)
	}

	if shuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Printf("Press 'Enter' to start...")
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	correct := evaluateProblems(problems, timer)
	fmt.Printf("Your score is %d out of %d.\n", correct, len(problems))
}
