package main

import (
	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/utils"
)

var (
	logger log.Logger
)

func main() {
	logger = utils.InitializeDefaultLogger()
	logger.Log("msg", "Starting trends server..")
}