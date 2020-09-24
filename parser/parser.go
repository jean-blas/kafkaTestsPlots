package parser

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TS struct {
	Ts1, Ts2 int64
}

// Parse the file
// extract the timestamps
func ParseData(filename string) ([]int64, []int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var ts1, ts2 []int64
	lines := make([]TS, 0)
	// Read the lines
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}

		fields := strings.Split(string(line), ";")
		if len(fields) != 3 {
			return nil, nil, errors.New("Bad formatted line : " + string(line))
		}
		fromTS, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		toTS, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		lines = append(lines, TS{Ts1: fromTS, Ts2: toTS})
	}
	// sort the data according to timestamp
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].Ts1 < lines[j].Ts2
	})
	// return the results
	for _, ts := range lines {
		ts1 = append(ts1, ts.Ts1)
		ts2 = append(ts2, ts.Ts2)
	}
	return ts1, ts2, nil
}

// Parse the file
// extract the timestamps
// return the slice of diffs
func ParseAndDiff(filename string) ([]int64, error) {
	ts1, ts2, err := ParseData(filename)
	if err != nil {
		return nil, err
	}

	var records []int64

	for i := range ts1 {
		diff := ts2[i] - ts1[i]
		records = append(records, diff)
	}
	return records, nil
}
