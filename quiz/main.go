package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "path to the CSV file with problems to solve")

	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit("could not open provided file: %v\n", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = 2

	scanner := bufio.NewScanner(os.Stdin)

	count := 0
	correct := 0

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

		if !scanner.Scan() {
			break
		}

		answer := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(answer, p.a) {
			correct++
		}
	}

	if err := scanner.Err(); err != nil {
		exit("could not read the standard input: %v\n", err)
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
