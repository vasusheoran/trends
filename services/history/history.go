package history

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/utils"
)

type History interface {
	Read(sasSymbol string) ([]contracts.Stock, error)
	Write(sasSymbol string, listings []contracts.Stock) error
	UploadFile(symbol string, r *http.Request) error
}

type historyDataIndex struct {
	HP   int
	LP   int
	CP   int
	Date int
}

type history struct {
	logger log.Logger
	dbSvc  database.Database
}


func (s *history) UploadFile(symbol string, r *http.Request) error {
	file, handler, err := r.FormFile("file")
	if err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return err
	}
	defer file.Close()

	level.Info(s.logger).
		Log("msg", "file uploaded successfully", "handler", handler.Filename, "symbol")
	f, err := os.OpenFile(utils.HistoricalFilePath(symbol), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return err
	}

	level.Info(s.logger).
		Log("msg", "file uploaded successfully", "handler_filename", handler.Filename, "path", utils.HistoricalFilePath(symbol))
	return nil
}

func (s *history) Read(sasSymbol string) ([]contracts.Stock, error) {
	data, err := s.dbSvc.Read(utils.HistoricalFilePath(sasSymbol))
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
		return nil, err
	}

	st := s.parseData(data)

	s.Write(sasSymbol+"-1", st)

	return st, nil
}

func (s *history) Write(sasSymbol string, st []contracts.Stock) error {
	if st == nil {
		return errors.New("failed to write nil history")
	}
	var data [][]string

	data = append(data, []string{"Date", "CP", "HP", "LP"})

	for _, val := range st {
		var temp []string

		temp = append(temp, val.Date)
		temp = append(temp, fmt.Sprintf("%v", val.CP))
		temp = append(temp, fmt.Sprintf("%v", val.HP))
		temp = append(temp, fmt.Sprintf("%v", val.LP))

		data = append(data, temp)
	}

	return s.dbSvc.Write(utils.HistoricalFilePath(sasSymbol), data)
}

func (s *history) parseHeaders(records [][]string, index *historyDataIndex) {
	if records == nil {
		return
	}

	for i, val := range records[0] {
		switch val {
		case "cp", "CP", "Close", "close":
			index.CP = i
		case "hp", "HP", "High", "high":
			index.HP = i
		case "lp", "LP", "Low", "low":
			index.LP = i
		case "Date", "date":
			index.Date = i
		default:
			level.Warn(s.logger).Log("msg", "Column not found or supported", "name", val)
		}
	}
}

func (s *history) parseDate(dateStr string) (time.Time, error) {
	// Define the layout of the date string
	layout := "01-Jan-06"

	// Parse the date string using the layout
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Return the parsed date
	return date, nil
}

func (s *history) parseData(records [][]string) []contracts.Stock {
	if records == nil {
		return []contracts.Stock{}
	}

	level.Info(s.logger).Log("msg", "parsing CSV data")

	var index historyDataIndex
	s.parseHeaders(records, &index)

	records = append(records[:0], records[1:]...)

	var data []contracts.Stock
	var err error

	for i, row := range records {
		if len(row) == 0 {
			level.Error(s.logger).Log("err", "no data found", "row", i)
			continue
		}
		var temp contracts.Stock

		if temp.Time, err = s.parseDate(row[index.Date]); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.CP, err = strconv.ParseFloat(row[index.CP], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.HP, err = strconv.ParseFloat(row[index.HP], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		if temp.LP, err = strconv.ParseFloat(row[index.LP], 64); err != nil {
			level.Error(s.logger).Log("err", err.Error(), "date", temp.Date)
			continue
		}
		temp.Date = row[index.Date]
		data = append(data, temp)
	}

	sort.Sort(contracts.ByTime(data))
	return data
}

func New(logger log.Logger, db database.Database) History {
	return &history{
		logger: logger,
		dbSvc:  db,
	}
}
