package database

import (
	"errors"
	"fmt"

	//"database/sql"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/pkg/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SQLDatastore struct {
	logger log.Logger
	db     *gorm.DB
}

func (s *SQLDatastore) DeleteStocks(ticker string) error {
	tx := s.db.Model(contracts.Stock{}).Begin()
	result := tx.Where("ticker = ?", ticker).Delete(&contracts.Stock{})
	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return result.Error
	}

	result = tx.Commit()
	if result.Error == nil {
		return nil
	}

	s.logger.Log("Error deleting stocks", result.Error, "ticker", ticker)
	result = tx.Rollback()
	if result.Error != nil {
		s.logger.Log("Error rolling transaction back", result.Error)
		return result.Error
	}

	return nil
}

func (s *SQLDatastore) GetDistinctTicker(pattern string) ([]string, error) {
	var stocks []contracts.Stock
	result := s.db.Model(contracts.Stock{}).Select("ticker").Where("ticker LIKE ?", "%"+pattern+"%").Distinct("ticker").Find(&stocks)
	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	var tickers []string
	for _, stock := range stocks {
		tickers = append(tickers, stock.Ticker)
	}

	return tickers, nil
}

func (s *SQLDatastore) ReadStockByTicker(ticker string) ([]contracts.Stock, error) {
	var stocks []contracts.Stock
	result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Order("time desc").Find(&stocks)
	//result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Limit(500).Find(&stocks)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	if len(stocks) == 0 {
		return nil, errors.New(fmt.Sprintf("failed to fetch stocks for `%s`", ticker))
	}

	return stocks, nil
}

func (s *SQLDatastore) SaveStocks(data []contracts.Stock) error {
	result := s.db.Model(contracts.Stock{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ticker"}, {Name: "date"}}, // key column
		UpdateAll: true,
		//DoUpdates: clause.AssignmentColumns([]string{"close", "high", "low", "time"}), // column needed to be updated
	}).Create(&data)
	if result.Error != nil {
		s.logger.Log("msg", "Error saving stocks", "error", result.Error)
		return result.Error
	}

	return nil

	//tx := s.db.Model(contracts.Stock{}).Begin()
	//result := tx.Clauses(clause.OnConflict{
	//	Columns:   []clause.Column{{Name: "ticker"}, {Name: "date"}},                  // key column
	//	DoUpdates: clause.AssignmentColumns([]string{"close", "high", "low", "time"}), // column needed to be updated
	//}).Create(data)
	//
	//if result.Error != nil {
	//	s.logger.Log("error", result.Error)
	//	return result.Error
	//}
	//
	//result = tx.Commit()
	//if result.Error == nil {
	//	return nil
	//}
	//
	//s.logger.Log("Error saving stocks", result.Error)
	//result = tx.Rollback()
	//if result.Error != nil {
	//	s.logger.Log("Error rolling transaction back", result.Error)
	//	return result.Error
	//}
	//
	//return nil
}

func (s *SQLDatastore) UpdateStock(data contracts.Stock) error {
	result := s.db.Save(&data)
	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return result.Error
	}

	s.db.Commit()

	return nil
}

func NewSqlDatastore(logger log.Logger, dbPath string) (*SQLDatastore, error) {
	if len(dbPath) == 0 {
		dbPath = "data/test.db"
	}

	err := os.MkdirAll(filepath.Dir(dbPath), 700)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&contracts.Summary{})

	db.AutoMigrate(&contracts.Stock{})

	return &SQLDatastore{logger: logger, db: db}, nil
}
