package main

import "testing"
import "reflect"

func Test_parseSort(t *testing.T) {
	tests := []struct {
		Input  string
		Output *sortOption
	}{
		{
			Input: "id desc",
			Output: &sortOption{
				column:   "id",
				sortType: "DESC",
			},
		},
		{
			Input: "fuga_id desc",
			Output: &sortOption{
				column:   "fuga_id",
				sortType: "DESC",
			},
		},
		{
			Input: "foo asc",
			Output: &sortOption{
				column:   "foo",
				sortType: "ASC",
			},
		},
		{
			Input:  "fuga_id",
			Output: nil,
		},
		{
			Input:  "fuga_id hoge",
			Output: nil,
		},
	}

	for _, test := range tests {
		if !reflect.DeepEqual(parseSort(test.Input), test.Output) {
			t.Fatal("error invalid result")
		}
	}
}
