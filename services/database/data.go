package database

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Database interface {
	Read(path string) ([][]string, error)
	Write(path string, data [][]string) error
}

type DB struct {
	logger log.Logger
}

func (d *DB) Read(file string) ([][]string, error) {
	_, err := os.Stat(file)
	if err != nil {
		level.Error(d.logger).Log("msg", "file does not exist", "err", err.Error())
		return nil, err
	}
	return d.csvReader(file)
}

func (d *DB) Write(path string, data [][]string) error {
	err := createFile(path)
	if err != nil {
		level.Error(d.logger).Log("msg", "Failed to create/open input file "+path, "err", err)
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		level.Error(d.logger).Log("msg", "Unable to read input file "+path, "err", err)
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	err = csvWriter.WriteAll(data)
	if err != nil {
		level.Error(d.logger).Log("msg", "Unable to write to csv file for "+path, "err", err)
		return err
	}

	return nil
}

func (d *DB) csvReader(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		level.Error(d.logger).Log("msg", "Unable to read input file "+path, "err", err)
		return nil, err
	}
	defer f.Close()
	records := [][]string{}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			// Reached end of file, no more lines
			break
		} else if errors.Is(err, csv.ErrFieldCount) {
			level.Error(d.logger).Log("msg", "skipping record", "err", err.Error())
			continue
		} else if err != nil {
			// Handle any other error
			level.Error(d.logger).Log("msg", "failed to read record", "err", err.Error())
			return nil, err
		}
		// Check if line is empty (only whitespace and newline)
		isEmpty := true
		for _, field := range record {
			if len(strings.TrimSpace(field)) > 0 {
				isEmpty = false
				break
			}
		}

		if !isEmpty {
			// Not an empty line, return the record
			records = append(records, record)
		}
	}

	return records, nil
}

func NewDatabase(logger log.Logger) Database {
	return &DB{logger: logger}
}

func createFile(path string) error {
	// check if file exists
	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}
