package listing

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/utils"
)

type listingsIndex struct {
	Name   int
	SAS    int
	Symbol int
	Series int
}

type listing struct {
	logger log.Logger
	dbSvc  api.Database
}

func (s *listing) Read() []contracts.Listing {
	data, err := s.dbSvc.Read(utils.SymbolsFilePath())
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
	}

	return s.parseData(data)
}

func (s *listing) Write(listings []contracts.Listing) error {
	var data [][]string

	data = append(data, []string{"Name", "SAS", "Symbol", "Series"})

	for _, val := range listings {
		var temp []string

		temp = append(temp, val.Name)
		temp = append(temp, val.SASSymbol)
		temp = append(temp, val.Symbol)
		temp = append(temp, val.Series)

		data = append(data, temp)
	}

	return s.dbSvc.Write(constants.SymbolsFilePath, data)
}

func (s *listing) parseData(records [][]string) []contracts.Listing {
	if records == nil {
		return []contracts.Listing{}
	}

	var index listingsIndex
	s.parseHeaders(records, &index)

	records = append(records[:0], records[1:]...)

	var data []contracts.Listing
	for _, row := range records {
		var temp contracts.Listing

		temp.Name = row[index.Name]
		temp.Series = row[index.Series]
		temp.SASSymbol = row[index.SAS]
		temp.Symbol = row[index.Symbol]

		data = append(data, temp)
	}

	return data
}

func (s *listing) parseHeaders(records [][]string, index *listingsIndex) {
	if records == nil {
		return
	}
	for i, val := range records[0] {
		switch val {
		case "Name", "name", "Company", "company":
			index.Name = i
		case "SAS", "sas", "sas_symbol":
			index.SAS = i
		case "symbol", "Symbol":
			index.Symbol = i
		case "Series", "series":
			index.Series = i
		default:
			level.Warn(s.logger).Log("msg", "Column not found or supported", "name", val)
		}
	}
}

func New(logger log.Logger, db api.Database) api.ListingsAPI {
	return &listing{
		logger: logger,
		dbSvc:  db,
	}
}
