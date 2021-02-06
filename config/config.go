package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BindAddr string
	LogLevel log.Level
	Database *Database
	RedisUrl string
	Jwt      *Jwt
}

func NewConfig() *Config {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal("Incorrect log level")
	}
	log.SetLevel(logLevel)

	return &Config{
		BindAddr: os.Getenv("BIND_ADDR"),
		LogLevel: logLevel,
		Database: newDatabase(),
		RedisUrl: os.Getenv("REDIS_URL"),
		Jwt:      newJwt(),
	}
}

type Database struct {
	Host            string
	Port            uint16
	Name            string
	User            string
	Password        string
	SslMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration
}

func newDatabase() *Database {
	port, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		log.Fatal("Incorrect database port")
	}

	maxOpenConns, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_OPEN_CONNS"), 10, 32)
	if err != nil {
		log.Fatal("Incorrect database max open connections")
	}

	maxIdleConns, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_IDLE_CONNS"), 10, 32)
	if err != nil {
		log.Fatal("Incorrect database max idle connections")
	}

	connMaxLifeTime, err := time.ParseDuration(os.Getenv("DATABASE_CONN_MAX_LIFE_TIME"))
	if err != nil {
		log.Fatal("Incorrect database connection max life time")
	}

	connMaxIdleTime, err := time.ParseDuration(os.Getenv("DATABASE_CONN_MAX_IDLE_TIME"))
	if err != nil {
		log.Fatal("Incorrect database connection max idle time")
	}

	return &Database{
		Host:            os.Getenv("DATABASE_HOST"),
		Port:            uint16(port),
		Name:            os.Getenv("DATABASE_NAME"),
		User:            os.Getenv("DATABASE_USER"),
		Password:        os.Getenv("DATABASE_PASSWORD"),
		SslMode:         os.Getenv("DATABASE_SSL_MODE"),
		MaxOpenConns:    int(maxOpenConns),
		MaxIdleConns:    int(maxIdleConns),
		ConnMaxIdleTime: connMaxIdleTime,
		ConnMaxLifeTime: connMaxLifeTime,
	}
}

type Jwt struct {
	AccessSecret  string
	RefreshSecret string
}

func newJwt() *Jwt {
	return &Jwt{
		AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
	}
}
