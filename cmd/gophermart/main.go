package main

import (
	"fmt"
	"github.com/vladislaoramos/gophermart/configs"
	"github.com/vladislaoramos/gophermart/internal/app"
	"github.com/vladislaoramos/gophermart/pkg/logger"
	stdLog "log"
	"os"
)

func main() {
	cfg := configs.NewConfig()

	f, err := os.OpenFile("/tmp/log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		stdLog.Fatal("unable to open file for log")
	}

	log := logger.New(cfg.Logger.Level, f)
	log.Info(fmt.Sprintf("%+v", cfg))

	app.Run(cfg, log)
}
