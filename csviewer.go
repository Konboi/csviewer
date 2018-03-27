package main

import (
	"log"
	"os"
	"sort"
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

type sortOption struct {
	column   string
	sortType string
}

func newSortData(data map[int]string, sortType string) *sortData {
	d := &sortData{
		data:     data,
		sortType: strings.ToUpper(sortType),
	}

	for k := range data {
		d.index = append(d.index, k)
	}

	sort.Ints(d.index)

	return d
}

type sortData struct {
	data     map[int]string
	index    []int
	sortType string
}

func (sd *sortData) Len() int {
	return len(sd.data)
}

func (sd *sortData) Less(i, j int) bool {
	iData := sd.data[i]
	jData := sd.data[j]

	iVal, errI := strconv.Atoi(iData)
	jVal, errJ := strconv.Atoi(jData)
	if errI != nil || errJ != nil {
		// iData and jData: string
		if errI != nil && errJ != nil {
			sorted := sort.StringsAreSorted([]string{iData, jData})

			if sd.sortType == "DESC" {
				sorted = !sorted
			}

			return sorted
		}
		// iData is number but jData is maybe empty
		if errI == nil && errJ != nil {
			jVal = 0
		}
		// jData is number but iData is maybe empty
		if errI != nil && errJ == nil {
			iVal = 0
		}
	}

	if sd.sortType == "DESC" {
		return iVal > jVal
	}

	return iVal < jVal // ASC
}

func (sd *sortData) Swap(i, j int) {
	sd.index[i], sd.index[j] = sd.index[j], sd.index[i]
	sd.data[i], sd.data[j] = sd.data[j], sd.data[i]
}

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
			log.Fatal("filter format is invalid: ", f)
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

func (v *Csviewer) Print(opt *sortOption) {
	var printRows [][]string
	sortIndexAndValue := make(map[int]string)
	var sortData *sortData

	count := 0
	for i, rowMap := range v.rowsMap {
		if v.filter(rowMap) {
			var row []string
			for _, j := range v.printInfo.printColumnIndex {
				row = append(row, v.rows[i][j])

				if opt != nil && opt.column == v.printInfo.printColumns[j] {
					sortIndexAndValue[i] = v.rows[i][j]
				}
			}
			printRows = append(printRows, row)
			count++
		}

		if 0 < v.limit && v.limit <= count {
			break
		}
	}

	if opt != nil {
		sortData = newSortData(sortIndexAndValue, opt.sortType)
		sort.Sort(sortData)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(v.printInfo.printColumns)
	if opt != nil {
		rows := [][]string{}
		for _, i := range sortData.index {
			rows = append(rows, printRows[i])
		}
		t.AppendBulk(rows)
	} else {
		t.AppendBulk(printRows)
	}
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
