package database

import (
	"encoding/csv"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/api"
)

type DB struct {
	logger log.Logger
}

func (d *DB) Read(file string) ([][]string, error) {
	var _, err = os.Stat(file)
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

	csvReader := csv.NewReader(f)
	return csvReader.ReadAll()
}

func NewDatabase(logger log.Logger) api.Database {
	return &DB{logger: logger}
}

func createFile(path string) error {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}
