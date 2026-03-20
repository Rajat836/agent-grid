package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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

type ServiceClientConfig struct {
	ProtocolType string `yaml:"protocol"`
	HttpHost     string `yaml:"http_host"`
	GrpcHost     string `yaml:"grpc_host"`
}

type SQSConfig struct {
	Region string         `yaml:"region"`
	Events SQSEventConfig `yaml:"events"`
}

type EventConfig struct {
	Queue string `yaml:"queue"`
	Topic string `yaml:"topic"`
}

type SQSEventConfig struct {
	CustomEventCreation        EventConfig `yaml:"custom_data_creation"`
	EntityUserDataCreation     EventConfig `yaml:"entity_user_data_creation"`
	EntityBusinessDataCreation EventConfig `yaml:"entity_business_data_creation"`
	EntityDocumentCreation     EventConfig `yaml:"entity_document_creation"`
	EntityLocationCreation     EventConfig `yaml:"entity_location_creation"`
	EntityRiskScoreCreation    EventConfig `yaml:"entity_risk_score_creation"`
	EntityRelationshipCreation EventConfig `yaml:"entity_relationship_creation"`
}

type ClickHouseConfig struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"dsn"`
}

type FeatureFlags struct {
	EnableAuditLog bool `yaml:"EnableAuditLog"`
}

// Config represents an application configuration.
type Config struct {
	ServerPort      int                 `yaml:"ServerPort"`
	GrpcPort        int                 `yaml:"GrpcPort"`
	AppName         string              `yaml:"AppName"`
	AppVersion      string              `yaml:"AppVersion"`
	BaseUrl         string              `yaml:"BaseUrl"`
	Environment     string              `yaml:"Environment"`
	OtlpExporterUrl string              `yaml:"OtlpExporterUrl"`
	Redis           RedisConfig         `yaml:"Redis"`
	Database        DatabaseConfig      `yaml:"Database"`
	Notification    ServiceClientConfig `yaml:"Notification"`
	SQS             SQSConfig           `yaml:"SQS"`
	ClickHouse      ClickHouseConfig    `yaml:"ClickHouse"`
	FeatureFlags    FeatureFlags        `yaml:"FeatureFlags"`
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string) (*Config, error) {
	// default config
	var c Config

	// load from YAML config file
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// os.Stdout.Write(bytes)
	fmt.Printf("trying to load config from file: %s\n", file)

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	if len(c.Database.MasterDatabaseDsn) == 0 {
		return nil, errors.New("configuration was not picked up from the file")
	}

	fmt.Printf("loaded config from file: %s\n", file)
	return &c, err
}
