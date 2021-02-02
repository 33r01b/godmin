package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BindAddr string
	LogLevel string
	Database Database
}

func NewConfig() *Config {
	return &Config{
		BindAddr: os.Getenv("BIND_ADDR"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Database: newDatabase(),
	}
}

type Database struct {
	Host                  string
	Port                  int16
	Name                  string
	Password              string
	SslMode               string
	MaxConnections        int8
	AcquireTimeoutSeconds int64
}

func newDatabase() Database {
	databasePort, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		log.Fatal("Incorrect database port")
	}

	databaseMaxConnections, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_CONNECTIONS"), 10, 8)
	if err != nil {
		log.Fatal("Incorrect database max connections")
	}

	databaseAcquireTimeoutSeconds, err := strconv.ParseInt(os.Getenv("DATABASE_ACQUIRE_TIMEOUT_SECONDS"), 10, 64)
	if err != nil {
		log.Fatal("Incorrect database acquire timeout seconds")
	}

	return Database{
		Host:                  os.Getenv("DATABASE_HOST"),
		Port:                  int16(databasePort),
		Name:                  os.Getenv("DATABASE_NAME"),
		Password:              os.Getenv("DATABASE_PASSWORD"),
		SslMode:               os.Getenv("DATABASE_SSL_MODE"),
		MaxConnections:        int8(databaseMaxConnections),
		AcquireTimeoutSeconds: databaseAcquireTimeoutSeconds,
	}
}
