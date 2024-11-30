package ticker

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/database"
	"github.com/vsheoran/trends/services/history"
	"github.com/vsheoran/trends/utils"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
)

var (
	logger log.Logger
)

const (
	resultDir = "test/result"
	testDir   = "test/input"
)

func iterateFiles(t *testing.T, logger log.Logger, dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		logger.Log("filename", filename)

		resultFile := filepath.Join(resultDir, strings.TrimSuffix(filename, ".csv")) + ".json"

		_, err = os.Stat(resultFile)
		if err != nil {
			return err
		}

		var summary contracts.Summary
		err = utils.ReadFromFile(logger, resultFile, &summary)
		if err != nil {
			return err
		}

		dbPath := "test/test.db"
		defer os.Remove(dbPath)
		sqlDB, err := database.NewSqlDatastore(logger, dbPath)
		if err != nil {
			t.Fatal(err)
		}

		db := database.NewCSVDatastore(logger)
		hs := history.New(logger, db, sqlDB)

		cardsSvc := cards.New(logger)
		ticker := NewTicker(logger, cardsSvc, hs)
		s, err := ticker.Init(filename, path)

		assert.Nil(t, err)
		assert.NotNil(t, s)

		//ticker.Update(filename, exp.In)

		//res, _ := ticker.Get(filename)

		assert.Equal(t, summary, s)

		return nil
	})

	return err
}

func assertFields(t *testing.T, actualSummary, expectedSummary contracts.Summary, validateCols []string) {
	vActualSummary := reflect.ValueOf(actualSummary)
	vExpectedSummary := reflect.ValueOf(expectedSummary)
	for _, col := range validateCols {
		expectedfield := vExpectedSummary.FieldByName(col)
		if !expectedfield.IsZero() {
			actualField := vActualSummary.FieldByName(col)
			assert.Equal(t, expectedfield.Equal(actualField), true)
		}
	}
}

func TestTicker_Init(t *testing.T) {
	logger = utils.InitializeDefaultLogger()

	err := iterateFiles(t, logger, testDir)
	if err != nil {
		t.Fatal(err)
	}

}
