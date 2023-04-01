package main

import (
	"fmt"
	"github.com/vladislaoramos/gophemart/configs"
	"github.com/vladislaoramos/gophemart/internal/app"
	"github.com/vladislaoramos/gophemart/pkg/logger"
	"log"
	"os"
)

func main() {
	cfg := configs.NewConfig()

	f, err := os.OpenFile("/tmp/log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("unable to open file for log")
	}

	l := logger.New(cfg.Logger.Level, f)
	l.Info(fmt.Sprintf("%+v", cfg))

	app.Run(cfg, l)
}
