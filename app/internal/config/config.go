package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	ServerPort      int                `yaml:"ServerPort"`
	GrpcPort        int                `yaml:"GrpcPort"`
	AppName         string             `yaml:"AppName"`
	AppVersion      string             `yaml:"AppVersion"`
	BaseUrl         string             `yaml:"BaseUrl"`
	OtlpExporterUrl string             `yaml:"OtlpExporterUrl"`
	Environment     string             `yaml:"Environment"`
	Redis           RedisConfig        `yaml:"Redis"`
	Database        DatabaseConfig     `yaml:"Database"`
	ClickHouse      ClickHouseConfig   `yaml:"ClickHouse"`
	Notification    NotificationConfig `yaml:"Notification"`
	FeatureFlag     FeatureFlagConfig  `yaml:"FeatureFlag"`

	Server ServerConfig `yaml:"-"` // Computed
}

type ServerConfig struct {
	HTTPPort int
	GRPCPort int
}

type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	TlsEnabled bool   `yaml:"tls_enabled"`
}

type DatabaseConfig struct {
	MasterDatabaseDsn string `yaml:"master_database_dsn"`
	SlaveDatabaseDsn  string `yaml:"slave_database_dsn"`
}

type ClickHouseConfig struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"dsn"`
}

type NotificationConfig struct {
	Protocol string `yaml:"protocol"`
	HttpHost string `yaml:"http_host"`
	GrpcHost string `yaml:"grpc_host"`
}

type FeatureFlagConfig struct {
	AuditLog bool `yaml:"audit_log"`
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config/local.yml"
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not read config file %s, using defaults: %v\n", configFile, err)
		return getDefaultConfig(), nil
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set computed values
	config.Server = ServerConfig{
		HTTPPort: config.ServerPort,
		GRPCPort: config.GrpcPort,
	}

	return config, nil
}

// getDefaultConfig returns default configuration for local development
func getDefaultConfig() *Config {
	return &Config{
		ServerPort:      4441,
		GrpcPort:        4442,
		AppName:         "ontology_bot",
		AppVersion:      "0.0.1",
		BaseUrl:         "http://localhost:4441",
		OtlpExporterUrl: "localhost:4318",
		Environment:     "local",
		Redis: RedisConfig{
			Host:       "localhost",
			Port:       6379,
			TlsEnabled: false,
		},
		Database: DatabaseConfig{
			MasterDatabaseDsn: "postgres://user:password@localhost:5432/ontology_bot?sslmode=disable",
			SlaveDatabaseDsn:  "postgres://user:password@localhost:5432/ontology_bot?sslmode=disable",
		},
		ClickHouse: ClickHouseConfig{
			Enabled: false,
			DSN:     "",
		},
		Notification: NotificationConfig{
			Protocol: "grpc",
			HttpHost: "http://localhost:4701",
			GrpcHost: "localhost:4702",
		},
		FeatureFlag: FeatureFlagConfig{
			AuditLog: true,
		},
		Server: ServerConfig{
			HTTPPort: 4441,
			GRPCPort: 4442,
		},
	}
}
