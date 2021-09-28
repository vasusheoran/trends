package database

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
)

const (
	DateFormat = "2006-01-02"
	//DateFormat = "02/01/06"
)

type ColumnIndex struct {
	HP   int
	LP   int
	CP   int
	Date int
}

type DB struct {
	logger log.Logger
}

func (d *DB) readCSV(path string, colIndex *ColumnIndex) ([][]string, error) {

	f, err := os.Open(path)
	if err != nil {
		level.Error(d.logger).Log("msg", "Unable to read input file "+path, "err", err)
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		level.Error(d.logger).Log("msg", "Unable to parse file as CSV for "+path, "err", err)
		return nil, err
	}

	for i, val := range records[0] {
		switch val {
		case "cp", "CP", "Close":
			colIndex.CP = i
		case "hp", "HP", "High":
			colIndex.HP = i
		case "lp", "LP", "Low":
			colIndex.LP = i
		case "Date", "date":
			colIndex.Date = i
		default:
			level.Warn(d.logger).Log("msg", "Column not found or supported", "name", val)
		}
	}

	records = append(records[:0], records[1:]...)

	return records, nil
}

func (d *DB) Read(file string) ([]contracts.Candle, error) {
	var index ColumnIndex

	records, err := d.readCSV(file, &index)
	if err != nil {
		return nil, err
	}

	var data []contracts.Candle

	for _, row := range records {
		var temp contracts.Candle

		if temp.Date, err = time.Parse(DateFormat, row[index.Date]); err != nil {
			level.Error(d.logger).Log("err", err.Error())
			continue
		}
		if temp.CP, err = strconv.ParseFloat(row[index.CP], 64); err != nil {
			level.Error(d.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.HP, err = strconv.ParseFloat(row[index.HP], 64); err != nil {
			level.Error(d.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.LP, err = strconv.ParseFloat(row[index.LP], 64); err != nil {
			level.Error(d.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}

		data = append(data, temp)
	}

	return data, nil
}

func (d *DB) Write(path string, data []contracts.Candle) error {
	panic("implement me")
}

func NewDatabase(logger log.Logger) api.Database {
	return &DB{logger: logger}
}
