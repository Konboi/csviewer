package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/olekukonko/tablewriter"
)

type Csviewer struct {
	columns     []string
	rows        [][]string
	printInfo   printInfo
	rowsMap     []map[string]interface{}
	filters     []string
	limit       int
	isFiltersOr bool
}

type printInfo struct {
	printColumns     []string
	printColumnIndex []int
}

type funcFilter func(map[string]interface{}) bool

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

func newCsviwer(columns []string, rows [][]string, printColumns string, filters []string, limit int, isFiltersOr bool) *Csviewer {
	return &Csviewer{
		columns:     columns,
		rows:        rows,
		printInfo:   parsePrintColumns(columns, printColumns),
		rowsMap:     sliceToMap(columns, rows),
		filters:     filters,
		limit:       limit,
		isFiltersOr: isFiltersOr,
	}
}

func sliceToMap(columns []string, rows [][]string) []map[string]interface{} {
	// id,name,email
	// 1, foo, foo@email.com
	// 2, fuga, fuga@email.com
	// ↓
	// {"id": "1", "name", "foo", "email": "foo@email.com"}{"id": "2", "name", "fuga", "email": "fuga@email.com"}
	data := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		rowMap := make(map[string]interface{})
		var val interface{}
		for i, column := range columns {
			tmpval, err := strconv.ParseFloat(row[i], 64)
			if err == nil {
				val = tmpval
			} else {
				val = row[i]
			}

			rowMap[column] = val
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

func (v *Csviewer) Print(opt *sortOption) {
	var printRows [][]string
	sortIndexAndValue := make(map[int]string)
	var sortData *sortData

	count := 0
	for i, rowMap := range v.rowsMap {
		if v.filter(rowMap) {
			var row []string
			for _, j := range v.printInfo.printColumnIndex {
				row = append(row, fmt.Sprint(v.rows[i][j]))

				if opt != nil && opt.column == v.printInfo.printColumns[j] {
					sortIndexAndValue[i] = fmt.Sprint(v.rows[i][j])
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

func (v *Csviewer) filter(rowMap map[string]interface{}) bool {
	if len(v.filters) == 0 {
		return true
	}

	values := make(map[string]interface{}, 0)

	for key, val := range rowMap {
		values[key] = val
	}

	filters := make([]string, 0)

	for _, f := range v.filters {
		filters = append(filters, f)
	}

	op := " && "
	if v.isFiltersOr {
		op = " || "
	}

	expression, err := govaluate.NewEvaluableExpression(strings.Join(filters, op))
	if err != nil {
		log.Fatal(err)
	}
	result, err := expression.Evaluate(values)
	if err != nil {
		log.Fatal(err)
	}

	return result.(bool)
}
