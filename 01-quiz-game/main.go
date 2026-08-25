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

func parseArguments() (csv string, limit int, shuffle bool) {
	flag.StringVar(&csv, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle the quiz order")
	flag.Parse()

	return csv, limit, shuffle
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

func evaluateProblems(problems [][]string, c chan<- struct{}) {
	for i, problem := range problems {
		var answer string

		fmt.Printf("Problem #%d: %s = ", i+1, problem[0])
		fmt.Scan(&answer)

		answer = strings.TrimSpace(answer)

		if strings.EqualFold(answer, problem[1]) {
			c <- struct{}{}
		}
	}
	close(c)
}

func main() {
	csvFile, limit, shuffle := parseArguments()

	problems, err := readFile(csvFile)
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

	var correctAnswers int
	answerCh := make(chan struct{})
	timer := time.NewTimer(time.Duration(limit) * time.Second)

	go evaluateProblems(problems, answerCh)

loop:
	for {
		select {
		case _, ok := <-answerCh:
			if !ok {
				break loop
			}
			correctAnswers++
		case <-timer.C:
			fmt.Println("")
			close(answerCh)
		}
	}

	fmt.Printf("Your score is %d out of %d.\n", correctAnswers, len(problems))
}
