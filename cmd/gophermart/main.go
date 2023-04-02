package main

import (
	"fmt"
	"github.com/vladislaoramos/gophemart/configs"
	"github.com/vladislaoramos/gophemart/internal/app"
	"github.com/vladislaoramos/gophemart/pkg/logger"
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
