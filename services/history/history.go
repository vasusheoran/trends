package history

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/vsheoran/trends/services/ticker/cards/models"

	"github.com/vsheoran/trends/services/database"
	"gorm.io/gorm"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/contracts"
)

type History interface {
	Read(symbol string) ([]models.Ticker, error)
	Write(path string, tickers []models.Ticker) error
	ParseFile(path string, file multipart.File) ([]models.Ticker, error)
}

type historyDataIndex struct {
	Close int
	Open  int
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

func (s *history) ParseFile(symbol string, file multipart.File) ([]models.Ticker, error) {
	level.Debug(s.logger).Log("msg", "Parsing uploaded file", "symbol", symbol)

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	tickers := s.parseData(symbol, records)
	return tickers, nil
}

func (s *history) Read(symbol string) ([]models.Ticker, error) {
	return s.sqlDB.ReadTickers(symbol, "", "")
}

func (s *history) Write(path string, tickers []models.Ticker) error {
	if tickers == nil {
		return errors.New("failed to write nil history")
	}

	return s.sqlDB.SaveTickers(tickers)
}

func (s *history) parseHeaders(records [][]string, index *historyDataIndex) {
	if records == nil {
		return
	}

	for i, val := range records[0] {
		toLowerVal := strings.Trim(strings.ToLower(val), " ")
		switch toLowerVal {
		case "w", "close":
			index.Close = i
		case "x", "open":
			index.Open = i
		case "y", "high":
			index.High = i
		case "z", "low":
			index.Low = i
		case "date":
			index.Date = i
		default:
			level.Warn(s.logger).Log("msg", "Column not found or supported", "name", val)
		}
	}
}

func (s *history) parseData(symbol string, records [][]string) []models.Ticker {
	if records == nil {
		return []models.Ticker{}
	}

	var index historyDataIndex
	s.parseHeaders(records, &index)

	records = append(records[:0], records[1:]...)

	var data []models.Ticker
	var err error
	var t time.Time

	for _, row := range records {
		var temp models.Ticker

		temp.Date = row[index.Date]

		t, err = s.parseDate(temp.Date)
		if err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}

		temp.Time = t

		if temp.W, err = strconv.ParseFloat(row[index.Close], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.X, err = strconv.ParseFloat(row[index.Open], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.Y, err = strconv.ParseFloat(row[index.High], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.Z, err = strconv.ParseFloat(row[index.Low], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}

		temp.Name = symbol
		temp.ParsedDate = temp.Time.Format("02-Jan-06")

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
