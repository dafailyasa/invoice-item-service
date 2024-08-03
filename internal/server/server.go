package server

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dafailyasa/invoice-item-service/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Echo   *echo.Echo
	Cfg    config.Config
	Logger zerolog.Logger
	ES     *es.Client
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		DB:     lo.Must(config.NewDatabase(cfg.Database)),
		ES:     lo.Must(config.NewESClient(cfg.ElasticSearch)),
		Echo:   config.NewEchoServer(cfg),
		Cfg:    cfg,
		Logger: config.NewLogger(),
	}
}

func (s *Server) Run() error {
	s.MapApps()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	go func() {
		<-quit

		log.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		lo.Must(s.DB.DB()).Close()
		s.Echo.Shutdown(ctx)
	}()

	return s.Echo.Start(fmt.Sprintf(":%s", s.Cfg.Server.Port))
}
