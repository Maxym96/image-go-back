package config

import (
	"os"
	"time"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	JwtSecret           string
	JwtTTL              time.Duration
	RabbitHost          string
	RabbitUser          string
	RabbitPassword      string
}

func GetConfiguration() Configuration {
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migrations"
	}
	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}
	staticFilesLocation, set := os.LookupEnv("FILES_LOCATION")
	if !set {
		staticFilesLocation = "file_storage"
	}

	return Configuration{
		DatabaseName:        os.Getenv("DB_NAME"),
		DatabaseHost:        os.Getenv("DB_HOST"),
		DatabaseUser:        os.Getenv("DB_USER"),
		DatabasePassword:    os.Getenv("DB_PASSWORD"),
		MigrateToVersion:    migrateToVersion,
		MigrationLocation:   migrationLocation,
		FileStorageLocation: staticFilesLocation,
		JwtTTL:              72 * time.Hour,
		RabbitHost:          os.Getenv("RABBIT_HOST"),
		RabbitUser:          os.Getenv("RABBIT_USER"),
		RabbitPassword:      os.Getenv("RABBIT_PASSWORD"),
	}
}
