package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Csviewer struct {
	columns     []string
	rows        [][]string
	printInfo   printInfo
	rowsMap     []map[string]string
	filters     []string
	funcFilters map[string]funcFilter
	limit       int
}

type printInfo struct {
	printColumns     []string
	printColumnIndex []int
}

type funcFilter func(string) bool

func newCsviwer(columns []string, rows [][]string, printColumns string, filters []string, limit int) *Csviewer {
	return &Csviewer{
		columns:     columns,
		rows:        rows,
		printInfo:   parsePrintColumns(columns, printColumns),
		rowsMap:     sliceToMap(columns, rows),
		filters:     filters,
		funcFilters: parseFilters(filters),
		limit:       limit,
	}
}

func sliceToMap(columns []string, rows [][]string) []map[string]string {
	// id,name,email
	// 1, foo, foo@email.com
	// 2, fuga, fuga@email.com
	// ↓
	// {"id": "1", "name", "foo", "email": "foo@email.com"}{"id": "2", "name", "fuga", "email": "fuga@email.com"}
	data := make([]map[string]string, 0, len(rows))
	for _, row := range rows {
		rowMap := make(map[string]string)
		for i, column := range columns {
			rowMap[column] = row[i]
		}
		data = append(data, rowMap)
	}

	return data
}

func parsePrintColumns(columns []string, showColumns string) printInfo {
	index := make([]int, 0)
	prints := make([]string, 0)

	if showColumns == "" {
		for i, _ := range columns {
			index = append(index, i)
		}
		return printInfo{columns, index}
	}

	for _, c := range strings.Split(strings.TrimSpace(showColumns), ",") {
		for i, column := range columns {
			if c == column {
				prints = append(prints, column)
				index = append(index, i)
				break
			}
		}
	}

	return printInfo{prints, index}
}

func parseFilters(filters []string) map[string]funcFilter {
	funcFilters := map[string]funcFilter{}
	for _, f := range filters {
		token := strings.SplitN(f, " ", 3)
		if len(token) < 3 {
			log.Fatalf("filter format is invalid", f)
		}
		key := token[0]
		switch token[1] {
		case ">", ">=", "<=", "<":
			r, err := strconv.ParseFloat(token[2], 64)
			if err != nil {
				log.Fatalf("error filter value column:'%s' check value:'%s'\n", key, token[2])
			}

			funcFilters[key] = func(val string) bool {
				if val == "" {
					return false
				}

				num, err := strconv.ParseFloat(val, 64)
				if err != nil {
					log.Printf("error check rows value:'%s'\n", val)
					return false
				}
				switch token[1] {
				case ">":
					return num > r
				case ">=":
					return num >= r
				case "<=":
					return num <= r
				case "<":
					return num < r
				default:
					return false
				}
			}
		case "==":
			funcFilters[key] = func(val string) bool {
				return val == token[2]
			}
		case "!=":
			funcFilters[key] = func(val string) bool {
				return val != token[2]
			}
		}
	}
	return funcFilters
}

func (v *Csviewer) Print() {
	var printRows [][]string

	count := 0
	for i, rowMap := range v.rowsMap {
		if v.filter(rowMap) {
			var row []string
			for _, j := range v.printInfo.printColumnIndex {
				row = append(row, v.rows[i][j])
			}
			printRows = append(printRows, row)
			count++
		}

		if 0 < v.limit && v.limit <= count {
			break
		}
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(v.printInfo.printColumns)
	t.AppendBulk(printRows)
	t.Render()
}

func (v *Csviewer) filter(rowMap map[string]string) bool {

	print := true
	for key, funcFilter := range v.funcFilters {
		if !funcFilter(rowMap[key]) {
			print = false
			break
		}
	}

	return print
}
