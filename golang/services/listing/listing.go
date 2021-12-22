package listing

import (
	"errors"
	"os"

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

func (s *listing) Patch(sasSymbol string, listing contracts.Listing) error {
	data := s.readData()

	if data == nil {
		return errors.New("failed to update listing")
	}

	var index listingsIndex
	var isListingUpdated bool
	s.parseHeaders(data, &index)

	for _, row := range data {
		if row[index.SAS] == sasSymbol {
			row[index.Name] = listing.Name
			row[index.Series] = listing.Series
			row[index.SAS] = listing.SASSymbol
			row[index.Symbol] = listing.Symbol
			isListingUpdated = true
			break
		}
	}

	if isListingUpdated {
		return s.dbSvc.Write(utils.SymbolsFilePath(), data)
	}

	return nil
}

func (s *listing) Put(sasSymbol string, listing contracts.Listing) error {
	data := s.readData()

	if data == nil {
		return errors.New("failed to update listing")
	}

	var index listingsIndex
	var temp []string
	s.parseHeaders(data, &index)

	for _, row := range data {
		if row[index.SAS] == sasSymbol {
			return errors.New("symbol exists")
		}
	}

	for i := 0; i < 4; i++ {
		switch i {
		case index.SAS:
			temp = append(temp, listing.SASSymbol)
		case index.Symbol:
			temp = append(temp, listing.Symbol)
		case index.Series:
			temp = append(temp, listing.Series)
		case index.Name:
			temp = append(temp, listing.Name)
		}
	}

	data = append(data, temp)
	return s.dbSvc.Write(utils.SymbolsFilePath(), data)
}

func (s *listing) Delete(sasSymbol string) error {
	data := s.readData()

	if data == nil {
		return errors.New("failed to update listing")
	}

	var index listingsIndex
	var isListingDeleted bool
	s.parseHeaders(data, &index)

	for i, row := range data {
		if row[index.SAS] == sasSymbol {
			data = append(data[:i], data[i+1:]...)
			isListingDeleted = true
			break
		}
	}

	if isListingDeleted {
		os.Remove(utils.SymbolsFilePath())
		return s.dbSvc.Write(utils.SymbolsFilePath(), data)
	}

	return nil

}

func (s *listing) Read() []contracts.Listing {
	return s.parseData(s.readData())
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

func (s *listing) readData() [][]string {
	data, err := s.dbSvc.Read(utils.SymbolsFilePath())
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to retieve listings", "err", err.Error())
		return [][]string{}
	}
	return data
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
