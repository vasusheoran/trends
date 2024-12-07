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

func Test_Cards(t *testing.T) {

	logger = utils.InitializeDefaultLogger()

}

func TestTicker_Init(t *testing.T) {
	logger = utils.InitializeDefaultLogger()

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
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

		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}

		hs := history.New(logger, sqlDB)
		err = hs.UploadFile(path, f)
		if err != nil {
			t.Fatal(err)
		}

		cardsSvc := cards.New(logger)
		tckr := NewTicker(logger, cardsSvc, hs)
		s, err := tckr.Init(filename, path)

		assert.Nil(t, err)
		assert.NotNil(t, s)

		//ticker.Add(filename, exp.In)

		//res, _ := ticker.Get(filename)

		assert.Equal(t, summary, s)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

}
