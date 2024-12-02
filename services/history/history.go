package history

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/vsheoran/trends/services/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/contracts"
)

type History interface {
	Read(path string) ([]contracts.Stock, error)
	Write(path string, listings []contracts.Stock) error
	UploadFile(path string, r *http.Request) error
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
	dbSvc  database.DataStore
	sqlDB  *database.SQLDatastore
}

func (s *history) UploadFile(symbol string, r *http.Request) error {
	file, handler, err := r.FormFile("file")
	if err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return err
	}
	defer file.Close()

	level.Debug(s.logger).Log("msg", "Parsing uploaded file", "handler", handler.Filename, "path")

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	stocks := s.parseData(symbol, records)
	return s.Write(symbol, stocks)
}

func (s *history) Read(path string) ([]contracts.Stock, error) {
	return s.sqlDB.ReadStockByTicker(path)
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
	//var data [][]string
	//
	//data = append(data, []string{"Date", "Close", "High", "Low"})
	//
	//for _, val := range st {
	//	var temp []string
	//
	//	temp = append(temp, val.Date)
	//	temp = append(temp, fmt.Sprintf("%v", val.Close))
	//	temp = append(temp, fmt.Sprintf("%v", val.High))
	//	temp = append(temp, fmt.Sprintf("%v", val.Low))
	//
	//	data = append(data, temp)
	//}
	//
	//return s.dbSvc.Write(path, data)
}

func (s *history) parseHeaders(records [][]string, index *historyDataIndex) {
	if records == nil {
		return
	}

	for i, val := range records[0] {
		switch strings.ToLower(val) {
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

func (s *history) parseDate(date string) (time.Time, error) {
	layout := "2-Jan-2006" // Reference layout for parsing

	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
		return t, err
	}

	return t, nil
}

func New(logger log.Logger, db database.DataStore, sqlDB *database.SQLDatastore) History {
	return &history{
		logger: logger,
		dbSvc:  db,
		sqlDB:  sqlDB,
	}
}
