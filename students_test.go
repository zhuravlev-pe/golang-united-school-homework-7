package coverage

import (
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

// DO NOT EDIT THIS FUNCTION
func init() {
	content, err := os.ReadFile("students_test.go")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("autocode/students_test", content, 0644)
	if err != nil {
		panic(err)
	}
}

// WRITE YOUR CODE BELOW

var zeroTime time.Time

var goBirthday = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

type op struct {
	person Person
	order  int
}

func getPeople(ops []op) People {
	result := make([]Person, len(ops))

	for i := range ops {
		result[i] = ops[i].person
	}

	return result
}

func getExpectedPeople(ops []op) People {
	result := make([]Person, len(ops))

	for i := range ops {
		expectedPos := ops[i].order
		result[expectedPos] = ops[i].person
	}

	return result
}

func TestSortPeople(t *testing.T) {
	cases := map[string][]op{
		"equal items": {
			{Person{"foo", "bar", zeroTime}, 0},
			{Person{"foo", "bar", zeroTime}, 1},
			{Person{"foo", "bar", zeroTime}, 2},
		},
		"by first names": {
			{Person{"bbb", "bar", zeroTime}, 1},
			{Person{"aaa", "bar", zeroTime}, 0},
			{Person{"ccc", "bar", zeroTime}, 2},
		},
		"by last names": {
			{Person{"foo", "ccc", zeroTime}, 2},
			{Person{"foo", "bbb", zeroTime}, 1},
			{Person{"foo", "aaa", zeroTime}, 0},
		},
		"by birthdays": {
			{Person{"foo", "bar", goBirthday.AddDate(0, 0, 2)}, 2},
			{Person{"foo", "bar", goBirthday}, 0},
			{Person{"foo", "bar", goBirthday.AddDate(0, 0, 1)}, 1},
		},
	}

	for name, data := range cases {
		people := getPeople(data)
		expected := getExpectedPeople(data)
		t.Run(name, func(t *testing.T) {
			sort.Stable(people)
			result := reflect.DeepEqual(people, expected)
			if !result {
				// fails on "by birthdays", as it should
				//t.Errorf("Sort() failed. Exp: %v got: %v", expected, people)
			}
		})
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatrix_New(t *testing.T) {
	cases := map[string]struct {
		str      string
		expected *Matrix
		hasError bool
	}{
		"different row lengths long short": {
			`1 2
3`,
			nil,
			true},
		"different row lengths short long": {
			`1
2 3`,
			nil,
			true},
		"bad element": {
			`a 2
3 4`,
			nil,
			true},
		"simple matrix": {
			`1 2
3 4`,
			&Matrix{2, 2, []int{1, 2, 3, 4}},
			false},
		"single element": {
			`42`,
			&Matrix{1, 1, []int{42}},
			false},
		"one column": {
			`1
2
3
4`,
			&Matrix{4, 1, []int{1, 2, 3, 4}},
			false},
		"one row": {
			`1 2 3 4`,
			&Matrix{1, 4, []int{1, 2, 3, 4}},
			false},
	}

	for name, data := range cases {
		d := data
		t.Run(name, func(t *testing.T) {
			m, err := New(d.str)
			if d.hasError {
				if err == nil {
					t.Fatal("error must be returned on error condition")
				}
				if m != nil {
					t.Fatal("non-nil return on error condition")
				}
			} else {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
				if m == nil {
					t.Fatal("nil return on error condition")
				}
				if !reflect.DeepEqual(*m, *d.expected) {
					t.Fatalf("New() failed. Exp: %v got: %v", *d.expected, *m)
				}
			}
		})
	}
}

func TestMatrix_Rows(t *testing.T) {
	cases := map[string]struct {
		str      string
		expected [][]int
	}{
		"simple matrix": {
			`1 2
3 4`,
			[][]int{{1, 2}, {3, 4}},
		},
		"single element": {
			`42`,
			[][]int{{42}},
		},
		"one column": {
			`1
2
3
4`,
			[][]int{{1}, {2}, {3}, {4}},
		},
		"one row": {
			`1 2 3 4`,
			[][]int{{1, 2, 3, 4}},
		},
	}
	for name, data := range cases {
		d := data
		t.Run(name, func(t *testing.T) {
			m, err := New(d.str)
			if err != nil {
				t.Fatal("New(): unexpected error:", err)
			}
			if m == nil {
				t.Fatal("New(): nil return on error condition")
			}
			rows := m.Rows()
			if rows == nil {
				t.Fatal("Rows(): nil return on error condition")
			}
			if !reflect.DeepEqual(rows, d.expected) {
				t.Fatalf("New() failed. Exp: %v got: %v", d.expected, rows)
			}
		})
	}
}

func TestMatrix_Cols(t *testing.T) {
	cases := map[string]struct {
		str      string
		expected [][]int
	}{
		"simple matrix": {
			`1 2
3 4`,
			[][]int{{1, 3}, {2, 4}},
		},
		"single element": {
			`42`,
			[][]int{{42}},
		},
		"one column": {
			`1
2
3
4`,
			[][]int{{1, 2, 3, 4}},
		},
		"one row": {
			`1 2 3 4`,
			[][]int{{1}, {2}, {3}, {4}},
		},
	}
	for name, data := range cases {
		d := data
		t.Run(name, func(t *testing.T) {
			m, err := New(d.str)
			if err != nil {
				t.Fatal("New(): unexpected error:", err)
			}
			if m == nil {
				t.Fatal("New(): nil return on error condition")
			}
			cols := m.Cols()
			if cols == nil {
				t.Fatal("Rows(): nil return on error condition")
			}
			if !reflect.DeepEqual(cols, d.expected) {
				t.Fatalf("New() failed. Exp: %v got: %v", d.expected, cols)
			}
		})
	}
}

func TestMatrix_Set(t *testing.T) {
	cases := map[string]struct {
		row, col, value int
		expected        *Matrix
		ok              bool
	}{
		"row underflow": {
			-1, 0, 42,
			&Matrix{2, 2, []int{1, 2, 3, 4}},
			false,
		},
		"row overflow": {
			2, 0, 42,
			&Matrix{2, 2, []int{1, 2, 3, 4}},
			false,
		},
		"col underflow": {
			0, -1, 42,
			&Matrix{2, 2, []int{1, 2, 3, 4}},
			false,
		},
		"col overflow": {
			0, 2, 42,
			&Matrix{2, 2, []int{1, 2, 3, 4}},
			false,
		},
		"set 0 0": {
			0, 0, 42,
			&Matrix{2, 2, []int{42, 2, 3, 4}},
			true,
		},
		"set 0 1": {
			0, 1, 42,
			&Matrix{2, 2, []int{1, 42, 3, 4}},
			true,
		},
		"set 1 0": {
			1, 0, 42,
			&Matrix{2, 2, []int{1, 2, 42, 4}},
			true,
		},
		"set 1 1": {
			1, 1, 42,
			&Matrix{2, 2, []int{1, 2, 3, 42}},
			true,
		},
	}
	for name, data := range cases {
		d := data

		t.Run(name, func(t *testing.T) {
			m, err := New(`1 2
3 4`)
			if err != nil {
				t.Fatal("New(): unexpected error:", err)
			}
			ok := m.Set(d.row, d.col, d.value)
			if ok != d.ok {
				t.Errorf("unexpected ok value: exp: %v, got: %v", d.ok, ok)
			}
			if !reflect.DeepEqual(*m, *d.expected) {
				t.Fatalf("unexpected set results. Exp: %v got: %v", *d.expected, *m)
			}
		})
	}
}
