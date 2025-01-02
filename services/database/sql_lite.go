package database

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vsheoran/trends/services/metrics"
	"strings"
	"time"

	"github.com/vsheoran/trends/services/ticker/cards/models"

	//"database/sql"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SQLDatastore struct {
	logger log.Logger
	db     *gorm.DB
	mr     *prometheus.Registry
}

type ORDER string

const (
	ORDER_DESC ORDER = "desc"
	LIMIT            = 10000
)

func (s *SQLDatastore) DeleteTicker(ticker string) error {
	startTime := time.Now()
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

	s.recordLatencyMetric(metrics.DeleteTickerLatency, startTime)
	return nil
}

func (s *SQLDatastore) GetDistinctTicker(pattern string) ([]string, error) {
	startTime := time.Now()
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

	s.recordLatencyMetric(metrics.GetUniqueTickerLatency, startTime)
	return tickers, nil
}

func (s *SQLDatastore) ReadTickers(ticker, pattern string, order ORDER) ([]models.Ticker, error) {
	startTime := time.Now()
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

	s.recordLatencyMetric(metrics.GetTickerByNameLatency, startTime)
	return tickers, nil
}

func (s *SQLDatastore) PaginateTickers(ticker, pattern string, offset, limit int, order ORDER) ([]models.Ticker, error) {
	startTime := time.Now()
	var tickers []models.Ticker

	db := s.db.Model(models.Ticker{}).
		Where("name = ?", ticker).
		Offset(offset).
		Order(fmt.Sprintf("time %s", order))

	if len(pattern) > 0 {
		db = db.Where("lower(date) LIKE ?", "%"+strings.ToLower(pattern)+"%")
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	result := db.Find(&tickers)

	if result.Error != nil {
		s.logger.Log("error", result.Error)
		return nil, result.Error
	}

	if len(tickers) == 0 {
		return nil, errors.New(fmt.Sprintf("failed to fetch stocks for `%s`", ticker))
	}

	s.recordLatencyMetric(metrics.PaginateTickerLatency, startTime)
	return tickers, nil
}

func (s *SQLDatastore) SaveTickers(data []models.Ticker) error {
	startTime := time.Now()
	result := s.db.Model(models.Ticker{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "parsed_date"}}, // key column
		UpdateAll: true,
		//DoUpdates: clause.AssignmentColumns([]string{"close", "high", "low", "time"}), // column needed to be updated
	}).Create(&data)
	if result.Error != nil {
		s.logger.Log("msg", "Error saving stocks", "error", result.Error)
		return result.Error
	}

	s.recordLatencyMetric(metrics.SaveTickerLatenct, startTime)
	return nil
}

func (s *SQLDatastore) recordLatencyMetric(name string, startTime time.Time) {
	metrics.GetSummary(
		name,
		"",
		s.mr,
		map[float64]float64{0.25: 0.1, 0.5: 0.1, 0.95: 0.1, 0.99: 0.1, 1.0: 0.1},
	).Observe(time.Since(startTime).Seconds())
}

func NewSqlDatastore(logger log.Logger, dbPath string, mr *prometheus.Registry) (*SQLDatastore, error) {
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

	return &SQLDatastore{logger: logger, db: db, mr: mr}, nil
}
