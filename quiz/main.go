package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	filePath := flag.String("f", "problems.csv", "path to the CSV file with problems to solve")

	flag.Parse()

	fmt.Printf("provided file is %q\n", *filePath)

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("could not open provided file: %v\n", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = 2

	scanner := bufio.NewScanner(os.Stdin)

	questionNumber := 0
	correctAnswers := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not read a record: %v\n", err)
		}

		questionNumber++

		question, answer := record[0], record[1]

		fmt.Printf("Problem #%d: %v = ", questionNumber, question)

		if !scanner.Scan() {
			break
		}

		if scanner.Text() == answer {
			correctAnswers++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("could not read console: %v\n", err)
	}

	fmt.Printf("You scored %d out of %d.\n", correctAnswers, questionNumber)
}
