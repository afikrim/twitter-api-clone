package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Host string `env:"APP_HOST" envDefault:"localhost"`
	Port int    `env:"APP_PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	DBHost        string `env:"DB_HOST" envDefault:"localhost"`
	DBPort        int    `env:"DB_PORT" envDefault:"3306"`
	DBUsername    string `env:"DB_USERNAME" envDefault:"root"`
	DBPassword    string `env:"DB_PASSWORD"`
	DBDatabase    string `env:"DB_DATABASE"`
	DBDebug       bool   `env:"DB_DEBUG" envDefault:"false"`
	DBDialect     string `env:"DB_DIALECT" envDefault:"mysql"`
	DBAutoMigrate bool   `env:"DB_AUTO_MIGRATE" envDefault:"true"`

	RedisHost      string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort      int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword  string `env:"REDIS_PASSWORD"`
	RedisDB        int    `env:"REDIS_DB" envDefault:"0"`
	RedisCacheDB   int    `env:"REDIS_CACHE_DB" envDefault:"1"`
	RedisSessionDB int    `env:"REDIS_SESSION_DB" envDefault:"2"`
}

func GetConfig() (*Config, error) {
	config := Config{}

	err := godotenv.Load(".env")
	err = env.Parse(&config)

	return &config, err
}
