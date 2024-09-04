package app

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// defines main session storage type based on server config given
type StorageType int

const (
	Memory StorageType = iota
	File
	Database
)

func (t StorageType) String() string {
	switch t {
	case Memory:
		return "Memory"
	case File:
		return "File"
	case Database:
		return "Database"
	}
	return fmt.Sprintf("Unknown (%d)", t)
}

type ServerConfig struct {
	Endpoint                 string
	MaxFileMemory            int64
	StorageMode              StorageType
	DatabaseDSN              string
	LogLevel                 string
	CompressibleContentTypes []string
	MaxConnectionRetries     uint64
}

func InitServerConfig() ServerConfig {

	cf := ServerConfig{}

	cf.CompressibleContentTypes = []string{
		"text/html",
		"application/json",
	}

	flag.StringVar(&cf.Endpoint, "a", "localhost:8080", "the address:port endpoint for server to listen")
	flag.Int64Var(&cf.MaxFileMemory, "m", 1, "Max memory size in MB to process files uploaded")
	flag.StringVar(&cf.DatabaseDSN, "d", "", "database DSN (format: 'host=<host> [port=port] user=<user> password=<xxxx> dbname=<mydb> sslmode=disable')")
	flag.StringVar(&cf.LogLevel, "l", "info", "log level")
	flag.Parse()

	if val, found := os.LookupEnv("ADDRESS"); found {
		cf.Endpoint = val
	}
	if val, found := os.LookupEnv("MAX_FILE_MEMORY"); found {
		intval, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			cf.MaxFileMemory = intval
		}
	}
	if val, found := os.LookupEnv("DATABASE_DSN"); found {
		cf.DatabaseDSN = val
	}
	if val, found := os.LookupEnv("LOG_LEVEL"); found {
		cf.LogLevel = val
	}

	if cf.Endpoint == "" {
		panic("PANIC: endpoint address:port is not set")
	}
	if cf.LogLevel == "" {
		panic("PANIC: log level is not set")
	}
	if cf.MaxFileMemory == 0 {
		panic("PANIC: Max file memory cannot be zero")
	}

	//set main storage type for current session
	cf.StorageMode = Database

	return cf
}
