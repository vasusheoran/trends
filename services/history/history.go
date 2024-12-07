package history

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/vsheoran/trends/services/database"
	"gorm.io/gorm"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/contracts"
)

type History interface {
	Read(path string) ([]contracts.Stock, error)
	Write(path string, listings []contracts.Stock) error
	UploadFile(path string, file multipart.File) error
}

type historyDataIndex struct {
	Close int
	High  int
	Low   int
	Date  int
}

type StocksORM struct {
	gorm.Model
	contracts.Stock
}

type history struct {
	logger log.Logger
	sqlDB  *database.SQLDatastore
}

func (s *history) UploadFile(symbol string, file multipart.File) error {
	level.Debug(s.logger).Log("msg", "Parsing uploaded file", "symbol", symbol)

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	stocks := s.parseData(symbol, records)
	return s.Write(symbol, stocks)
}

func (s *history) Read(path string) ([]contracts.Stock, error) {
	return s.sqlDB.ReadStockByTicker(path, "")
	//data, err := s.dbSvc.Read(path)
	//if err != nil {
	//	level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
	//	return nil, err
	//}
	//
	//return s.parseData(data), nil
}

func (s *history) Write(path string, stocks []contracts.Stock) error {
	if stocks == nil {
		return errors.New("failed to write nil history")
	}

	return s.sqlDB.SaveStocks(stocks)
}

func (s *history) parseHeaders(records [][]string, index *historyDataIndex) {
	if records == nil {
		return
	}

	for i, val := range records[0] {
		toLowerVal := strings.Trim(strings.ToLower(val), " ")
		switch toLowerVal {
		case "cp", "close":
			index.Close = i
		case "hp", "high":
			index.High = i
		case "lp", "low":
			index.Low = i
		case "date":
			index.Date = i
		default:
			level.Warn(s.logger).Log("msg", "Column not found or supported", "name", val)
		}
	}
}

func (s *history) parseData(symbol string, records [][]string) []contracts.Stock {
	if records == nil {
		return []contracts.Stock{}
	}

	var index historyDataIndex
	s.parseHeaders(records, &index)

	records = append(records[:0], records[1:]...)

	var data []contracts.Stock
	var err error
	var t time.Time

	for _, row := range records {
		var temp contracts.Stock

		temp.Date = row[index.Date]

		t, err = s.parseDate(temp.Date)
		if err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}

		temp.Time = t

		if temp.Close, err = strconv.ParseFloat(row[index.Close], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.High, err = strconv.ParseFloat(row[index.High], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.Low, err = strconv.ParseFloat(row[index.Low], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}

		temp.Ticker = symbol

		data = append(data, temp)
	}
	return data
}

func (s *history) parseDate(dateString string) (time.Time, error) {
	formats := []string{
		"2-Jan-2006",
		"02-Jan-2006",
		"2-Jan-06",
		"02-Jan-06",
		"2-01-2006",
		"02-1-2006",
		"2-01-06",
		"02-1-06",
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateString)
}

func New(logger log.Logger, sqlDB *database.SQLDatastore) History {
	return &history{
		logger: logger,
		sqlDB:  sqlDB,
	}
}
