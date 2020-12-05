package main

import (
	"bufio"
	"context"
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

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Press 'Enter' to start the quiz...")
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			exit("could not read the standard input: %v\n", err)
		}
	}

	answerCh := make(chan string)
	errCh := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count := 0
	correct := 0

	go scanAnswers(ctx, scanner, answerCh, errCh)

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

		select {
		case answer := <-answerCh:
			if strings.EqualFold(answer, p.a) {
				correct++
			}
		case err := <-errCh:
			exit("\ncould not read the standard input: %v\n", err)
		case <-ctx.Done():
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

type scanner interface {
	Scan() bool
	Text() string
	Err() error
}

func scanAnswers(ctx context.Context, s scanner, answerCh chan<- string, errCh chan<- error) {
	for s.Scan() {
		select {
		case answerCh <- strings.TrimSpace(s.Text()):
		case <-ctx.Done():
			return
		}
	}

	select {
	case errCh <- s.Err():
	case <-ctx.Done():
	}
}
