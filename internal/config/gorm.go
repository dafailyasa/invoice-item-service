package config

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewDatabase(cfg DatabaseConfig) (*gorm.DB, error) {

	otps := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
		TranslateError:       true,
		FullSaveAssociations: true,
	}

	if !cfg.Debug {
		otps.Logger = logger.Default.LogMode(logger.Silent)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FJakarta", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), otps)
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(cfg.Pool.Idle)
	connection.SetMaxOpenConns(cfg.Pool.Max)
	connection.SetConnMaxLifetime(time.Second * time.Duration(cfg.Pool.Lifetime))

	if !cfg.Debug {
		db.Logger.LogMode(logger.Silent)
	}

	return db, nil
}
