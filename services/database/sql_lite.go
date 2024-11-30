package database

import (
	//"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/pkg/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"path/filepath"
)

type SQLDatastore struct {
	logger log.Logger
	db     *gorm.DB
}

func (s SQLDatastore) ReadStockByTicker(ticker string) ([]contracts.Stock, error) {
	var stocks []contracts.Stock
	result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Order("created_at desc").Find(&stocks)
	//result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Limit(500).Find(&stocks)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	return stocks, nil
}

func (s SQLDatastore) SaveStocks(data []contracts.Stock) error {
	tx := s.db.Model(contracts.Stock{}).Begin()
	result := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ticker"}, {Name: "date"}},          // key column
		DoUpdates: clause.AssignmentColumns([]string{"close", "high", "low"}), // column needed to be updated
	}).Create(data)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return result.Error
	}

	result = tx.Commit()
	if result.Error == nil {
		return nil
	}

	s.logger.Log("Error saving stocks", result.Error)
	result = tx.Rollback()
	if result.Error != nil {
		s.logger.Log("Error rolling transaction back", result.Error)
		return result.Error
	}

	return nil
}

func (s SQLDatastore) UpdateStock(data contracts.Stock) error {
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
