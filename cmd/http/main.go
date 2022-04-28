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
	auth_service "github.com/afikrim/go-hexa-template/internal/core/services/auth"
	country_service "github.com/afikrim/go-hexa-template/internal/core/services/country"
	user_service "github.com/afikrim/go-hexa-template/internal/core/services/user"
	userfollowing_service "github.com/afikrim/go-hexa-template/internal/core/services/userfollowing"
	http_handler "github.com/afikrim/go-hexa-template/internal/handlers/http"
	country_repository "github.com/afikrim/go-hexa-template/internal/repositories/country"
	session_repository "github.com/afikrim/go-hexa-template/internal/repositories/session"
	user_repository "github.com/afikrim/go-hexa-template/internal/repositories/user"
	userfollowing_repository "github.com/afikrim/go-hexa-template/internal/repositories/userfollowing"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RedisConnType string

const (
	Default RedisConnType = "default"
	Cache                 = "cache"
	Session               = "session"
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

	redisSession, err := NewRedisInstance(cfg, Session)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.LstdFlags)

	countryRepository := country_repository.NewCountryRepository(db)
	sessionRepository := session_repository.NewSessionRepository(redisSession)
	userRepository := user_repository.NewUserRepository(db)
	userfollowingRepository := userfollowing_repository.NewUserFollowingRepository(db)

	authService := auth_service.NewAuthService(userRepository, sessionRepository)
	countryService := country_service.NewCountryService(countryRepository)
	userService := user_service.NewUserService(userRepository)
	userfollowingService := userfollowing_service.NewUserFollowingService(userfollowingRepository, userRepository)

	authHandler := http_handler.NewAuthHandler(authService)
	countryHandler := http_handler.NewCountryHandler(countryService)
	userHandler := http_handler.NewUserHandler(userService)
	userfollowingHandler := http_handler.NewUserFollowingHandler(userfollowingService)

	// Register routes
	apiV1Router := e.Group("/api/v1")
	authHandler.RegisterRoutes(apiV1Router)
	countryHandler.RegisterRoutes(apiV1Router)
	userHandler.RegisterRoutes(apiV1Router)
	userfollowingHandler.RegisterRoutes(apiV1Router)

	go func() {
		address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
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
		instance.AutoMigrate(&country_repository.Country{}, &user_repository.User{})
	}

	return instance, nil
}

func NewRedisInstance(config *config.Config, connType RedisConnType) (*redis.Client, error) {
	ctx := context.Background()
	redisConf := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	}
	switch connType {
	case Default:
		redisConf.DB = config.RedisDB
	case Cache:
		redisConf.DB = config.RedisCacheDB
	case Session:
		redisConf.DB = config.RedisSessionDB
	default:
		return nil, fmt.Errorf("unsupported redis connection type: %s", connType)
	}

	client := redis.NewClient(&redis.Options{
		DB: config.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
