package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/robfig/cron/v3"
)

// Processor processes template
type Processor struct {
	parser cron.Parser
}

// NewProcessor creates new processor
func NewProcessor() *Processor {
	return &Processor{}
}

func mustInclude(time time.Time, cronControl string) (bool, error) {
	p := cron.NewParser(cron.Dom | cron.Month | cron.Dow)
	sched, err := p.Parse(cronControl)
	if err != nil {
		return true, err
	}
	// cron schedule returns next instance, and since granularity is a day
	// we can check for yesterday
	next := sched.Next(time.AddDate(0, 0, -1))
	y1, m1, d1 := time.Date()
	y2, m2, d2 := next.Date()
	return y1 == y2 && m1 == m2 && d1 == d2, nil
}

// Process parses input from reader and writes the output to output for date
func (p *Processor) Process(date time.Time, reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(`<!--\s*@(.+)-->`)
	var err error
	var matches []string
	var ok bool
	for scanner.Scan() {
		matches = re.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 {
			ok, err = mustInclude(date, matches[1])
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
		}
		fmt.Fprintln(writer, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`NAME:
   daily-checklist - Daily checklist processor

USAGE:
   daily-checklist template_file_name [date]

   Processes the template and outputs the result to stdout. If date is not provided, current date will be used.

NOTES:
   Template file can be any text file. All lines will be printed to stdout by default.
   To control if a line should be outputted add "<!-- @ [dom] [mon] [dow] -->" segment to the line which will output it only if the date matches the template. Date matching uses cron formatting.`)
		os.Exit(0)
	}
	var err error

	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	date := time.Now()
	if len(os.Args) > 2 {
		layout := "2006-01-02"
		date, err = time.Parse(layout, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	}

	err = NewProcessor().Process(date, f, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
