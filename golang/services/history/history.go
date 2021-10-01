package history

import (
	"fmt"
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

func (s *history) Read(sasSymbol string) []contracts.Stock {
	data, err := s.dbSvc.Read(utils.HistoricalFilePath(sasSymbol))
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
	}

	return s.parseData(data)
}

func (s *history) Write(sasSymbol string, st []contracts.Stock) error {
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
