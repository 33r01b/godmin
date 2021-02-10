package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Port     uint16 `envconfig:"PORT" default:"8080" required:"true"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug" required:"true"`
	RedisUrl string `envconfig:"REDIS_URL" default:"localhost:6379" required:"true"`
	Database *Database
	Jwt      *Jwt
}

func NewConfig() *Config {
	conf := &Config{}

	err := envconfig.Process("", conf)
	if err != nil {
		log.Fatal("can't process the config: %w", err)
	}

	logLevel, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal("Incorrect log level")
	}
	log.SetLevel(logLevel)

	return conf
}

type Database struct {
	Host            string        `envconfig:"DATABASE_HOST" default:"localhost" required:"true"`
	Port            uint16        `envconfig:"DATABASE_PORT" default:"5432" required:"true"`
	Name            string        `envconfig:"DATABASE_NAME" default:"godmin_db_dev" required:"true"`
	User            string        `envconfig:"DATABASE_USER" default:"godmin" required:"true"`
	Password        string        `envconfig:"DATABASE_PASSWORD" default:"password" required:"true"`
	SslMode         string        `envconfig:"DATABASE_SSL_MODE" default:"disable" required:"true"`
	MaxOpenConns    int           `envconfig:"DATABASE_MAX_OPEN_CONNS" default:"1000" required:"true"`
	MaxIdleConns    int           `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"15" required:"true"`
	ConnMaxLifeTime time.Duration `envconfig:"DATABASE_CONN_MAX_IDLE_TIME" default:"30s" required:"true"`
	ConnMaxIdleTime time.Duration `envconfig:"DATABASE_CONN_MAX_LIFE_TIME" default:"0" required:"true"`
}

type Jwt struct {
	AccessSecret  string `envconfig:"JWT_ACCESS_SECRET" default:"secret;)" required:"true"`
	RefreshSecret string `envconfig:"JWT_REFRESH_SECRET" default:"secret;)" required:"true"`
}
