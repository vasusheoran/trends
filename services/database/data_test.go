package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/utils"
)

func TestDb_Read(t *testing.T) {
	logger := utils.InitializeDefaultLogger()
	service := csvDatastore{
		logger: logger,
	}
	const path = "/Users/vasusheoran/git/trends/golang/data/02-14-2021.csv"

	data, _ := service.Read(path)

	//level.Info(logger).Log(dateIndex, cpIndex, hpIndex, lpIndex)

	assert.NotNil(t, data)
}
