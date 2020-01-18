package main

import (
	"go.uber.org/zap"
	"os"
)

var (
	Revision string // Revision ... build時に注入する
	log      Logger
)

// Logger ... store zap logger
type Logger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

// logを初期化
func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	s := logger.Sugar()
	log = Logger{logger, s}
}

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		log.sugar.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
