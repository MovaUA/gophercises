package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "path to the CSV file with problems to solve")
	timeoutString := flag.String("limit", "30s", "time limit for the game")

	flag.Parse()

	timeout, err := time.ParseDuration(*timeoutString)
	if err != nil {
		exit("could not parse time limit: %v\n", err)
	}

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit("could not open provided file: %v\n", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = 2

	fmt.Print("Press 'Enter' to start the quiz...")
	fmt.Scanln()

	count := 0
	correct := 0

	timer := time.NewTimer(timeout)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			exit("could not read a record: %v\n", err)
		}

		count++

		p := problem{
			q: strings.TrimSpace(record[0]),
			a: strings.TrimSpace(record[1]),
		}

		fmt.Printf("Problem #%d: %s = ", count, p.q)

		answerCh := make(chan string)
		go getAnswer(answerCh)

		select {
		case answer := <-answerCh:
			if strings.EqualFold(answer, p.a) {
				correct++
			}
		case <-timer.C:
			fmt.Printf("\nTimeout of %v expired\n", timeout)
			for {
				_, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					exit("could not read a record: %v\n", err)
				}
				count++
			}
			break
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, count)
}

type problem struct {
	q string
	a string
}

func exit(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func getAnswer(answerCh chan<- string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	answerCh <- answer
}
