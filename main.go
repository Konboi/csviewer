package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/soh335/sliceflag"
)

func main() {
	var path, printColumns string
	var limit int
	var filters []string

	flag.StringVar(&path, "path", "", "set csv file path")
	flag.StringVar(&path, "p", "", "set csv file path")
	flag.IntVar(&limit, "limit", 0, "set max display rows  num")
	flag.IntVar(&limit, "l", 0, "set max display rows num")
	flag.StringVar(&printColumns, "columns", "", "print specify columns")
	flag.StringVar(&printColumns, "c", "", "print specify columns")
	sliceflag.StringVar(flag.CommandLine, &filters, "f", []string{}, "filter")
	sliceflag.StringVar(flag.CommandLine, &filters, "filter", []string{}, "filter")
	flag.Parse()

	d, err := loadData(path)
	if err != nil {
		log.Fatal("error load data", err.Error())
	}

	columns, rows, err := convertData(d)
	if err != nil {
		log.Fatal("error convet data", err.Error())
	}

	viewer := newCsviwer(columns, rows, printColumns, filters, limit)
	viewer.Print()
}

func loadData(path string) (io.Reader, error) {
	if path != "" {
		return os.Open(path)
	}

	stdin := os.Stdin
	stats, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if stats.Size() > 0 {
		return stdin, nil
	}

	return nil, fmt.Errorf("data is emptry")
}

func convertData(data io.Reader) ([]string, [][]string, error) {
	c := csv.NewReader(data)

	column, err := c.Read()
	if err != nil {
		return []string{}, [][]string{}, err
	}
	var rows [][]string
	for {
		row, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, [][]string{}, err
		}
		rows = append(rows, row)
	}

	return column, rows, nil
}
