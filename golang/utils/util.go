package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func InitializeDefaultLogger() log.Logger {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "scs_ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

func WriteToFile(logger log.Logger, data interface{}, filename string) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	err := ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
	}
}

func ReadFromFile(logger log.Logger, filename string, res interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
	}
}
