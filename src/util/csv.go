package util

import (
	"bufio"
	"encoding/csv"
	"heshi/errors"
	"io"
	"os"
)

func ParseCSVToMap(file string) (map[string][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a new reader.
	count := 0
	headers := []string{}
	headerColumns := make(map[string][]string)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		if count == 0 {
			headers = record
		} else {
			if len(headers) <= len(record) {
				for index := 0; index < len(headers); index++ {
					mapKey := headers[index]
					columnValues := headerColumns[mapKey]
					headerColumns[mapKey] = append(columnValues, record[index])
				}
			}
		}
		count++
	}
	return headerColumns, nil
}

type Row struct {
	Number  int
	Value   []string
	Message []string
	Ignored bool
}

func ParseCSVToStruct(file string) ([]Row, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a new reader.
	rs := []Row{}
	i := 0
	reader := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := reader.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		r := Row{
			Number:  i,
			Value:   record,
			Ignored: false,
		}
		rs = append(rs, r)
		i++
	}
	return rs, nil
}

func ParseCSVToArrays(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a new reader.
	rows := [][]string{}
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		rows = append(rows, record)
	}
	return rows, nil
}

func GetCSVHeaders(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		return record, nil
	}
	return nil, errors.Newf("fail to find headers in uploaded csv file %s", file)
}
