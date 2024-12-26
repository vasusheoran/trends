package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vsheoran/trends/services/ticker/cards/models"

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

type ORDER string

const (
	ORDER_DESC ORDER = "desc"
	LIMIT            = 10000
)

func (s *SQLDatastore) DeleteTicker(ticker string) error {
	tx := s.db.Model(models.Ticker{}).Begin()
	result := tx.Where("name = ?", ticker).Delete(&models.Ticker{})
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
	var stocks []models.Ticker
	result := s.db.Model(models.Ticker{}).Select("name").Where("name LIKE ?", "%"+pattern+"%").Distinct("name").Find(&stocks)
	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	var tickers []string
	for _, stock := range stocks {
		tickers = append(tickers, stock.Name)
	}

	return tickers, nil
}

func (s *SQLDatastore) ReadStockByTicker(ticker string, order ORDER) ([]contracts.Stock, error) {
	var stocks []contracts.Stock

	result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Order(fmt.Sprintf("time %s", order)).Find(&stocks)
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

func (s *SQLDatastore) ReadTickers(ticker, pattern string, order ORDER) ([]models.Ticker, error) {
	var tickers []models.Ticker

	result := s.db.Model(models.Ticker{}).Where("name = ?", ticker).Where("lower(date) LIKE ?", "%"+strings.ToLower(pattern)+"%").Order(fmt.Sprintf("time %s", order)).Find(&tickers)
	//result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Limit(500).Find(&stocks)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	if len(tickers) == 0 {
		return nil, errors.New(fmt.Sprintf("failed to fetch stocks for `%s`", ticker))
	}

	return tickers, nil
}

func (s *SQLDatastore) PaginateTickers(ticker, pattern string, offset, limit int, order ORDER) ([]models.Ticker, error) {

	var tickers []models.Ticker

	result := s.db.Model(models.Ticker{}).
		Where("name = ?", ticker).
		//Where("lower(date) LIKE ?", "%"+strings.ToLower(pattern)+"%").
		Offset(offset).
		Limit(limit).
		Order(fmt.Sprintf("time %s", order)).
		Find(&tickers)
	//result := s.db.Model(contracts.Stock{}).Where("ticker = ?", ticker).Limit(500).Find(&stocks)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	if len(tickers) == 0 {
		return nil, errors.New(fmt.Sprintf("failed to fetch stocks for `%s`", ticker))
	}

	return tickers, nil
}

func (s *SQLDatastore) SaveTickers(data []models.Ticker) error {
	result := s.db.Model(models.Ticker{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "parsed_date"}}, // key column
		UpdateAll: true,
		//DoUpdates: clause.AssignmentColumns([]string{"close", "high", "low", "time"}), // column needed to be updated
	}).Create(&data)
	if result.Error != nil {
		s.logger.Log("msg", "Error saving stocks", "error", result.Error)
		return result.Error
	}

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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		CreateBatchSize: 500,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Ticker{})

	//db.AutoMigrate(&contracts.Stock{})

	return &SQLDatastore{logger: logger, db: db}, nil
}
