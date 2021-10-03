package history

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

type historyDataIndex struct {
	HP   int
	LP   int
	CP   int
	Date int
}

type history struct {
	logger log.Logger
	dbSvc  api.Database
}

func (s *history) UploadFile(symbol string, r *http.Request) error {
	file, handler, err := r.FormFile("file_name")
	fileName := r.FormValue("file")
	if err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return err
	}
	defer file.Close()

	level.Info(s.logger).Log("msg", "file uploaded successfully", "handler", handler.Filename, "symbol")
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

	level.Info(s.logger).Log("msg", "file uploaded successfully", "name", fileName)
	return nil
}

func (s *history) Read(sasSymbol string) ([]contracts.Stock, error) {
	data, err := s.dbSvc.Read(utils.HistoricalFilePath(sasSymbol))
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
		return nil, err
	}

	return s.parseData(data), nil
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

func (s *history) parseData(records [][]string) []contracts.Stock {
	if records == nil {
		return []contracts.Stock{}
	}

	var index historyDataIndex
	s.parseHeaders(records, &index)

	records = append(records[:0], records[1:]...)

	var data []contracts.Stock
	var err error

	for _, row := range records {
		var temp contracts.Stock

		temp.Date = row[index.Date]
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

		data = append(data, temp)
	}
	return data
}

func New(logger log.Logger, db api.Database) api.HistoryAPI {
	return &history{
		logger: logger,
		dbSvc:  db,
	}
}
