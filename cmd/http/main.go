package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afikrim/go-hexa-template/config"
	country_service "github.com/afikrim/go-hexa-template/internal/core/services/country"
	http_handler "github.com/afikrim/go-hexa-template/internal/handlers/http"
	country_repository "github.com/afikrim/go-hexa-template/internal/repositories/country"
	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	db, err := NewDatabaseInstance(cfg)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	e := echo.New()
	e.Logger.SetLevel(log.LstdFlags)

	apiV1Router := e.Group("/api/v1")

	countryRepository := country_repository.NewCountryRepository(db)
	countryService := country_service.NewCountryService(countryRepository)
	countryHandler := http_handler.NewCountryHandler(countryService)

	// Register routes
	countryHandler.RegisterRoutes(apiV1Router)

	go func() {
		if err := e.Start(address); err != nil {
			log.Fatalf("Server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}

func NewDatabaseInstance(config *config.Config) (*gorm.DB, error) {
	var connectionString string
	gormConf := gorm.Config{}

	if config.DBDebug {
		gormConf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	var instance *gorm.DB
	var err error

	switch config.DBDialect {
	case "mysql":
		connectionString = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			config.DBUsername,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBDatabase,
		)
		instance, err = gorm.Open(mysql.Open(connectionString), &gormConf)
	case "postgres":
		connectionString = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			config.DBHost,
			config.DBPort,
			config.DBUsername,
			config.DBDatabase,
			config.DBPassword,
		)
		instance, err = gorm.Open(postgres.Open(connectionString), &gormConf)
	default:
		err = fmt.Errorf("unsupported database dialect: %s", config.DBDialect)
	}

	if err != nil {
		return nil, err
	}

	if config.DBDebug {
		return instance.Debug(), nil
	}

	if config.DBAutoMigrate {
		instance.AutoMigrate(&country_repository.Country{})
	}

	return instance, nil
}
