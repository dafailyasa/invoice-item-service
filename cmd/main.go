package main

import (
	"time"

	"github.com/dafailyasa/invoice-item-service/internal/config"
	"github.com/dafailyasa/invoice-item-service/internal/server"
	"github.com/samber/lo"
)

func main() {
	cfg := lo.Must(config.LoadConfigPath("config"))
	_ = lo.Must(time.LoadLocation(cfg.Server.TimeZone))

	server.NewServer(cfg).Run()
}
