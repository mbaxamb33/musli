package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	DBTestSource  string `mapstructure:"DB_TEST_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// It's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}

// GetDBSource returns the database source string based on the environment
func GetDBSource(isTest bool) string {
	// Check if we're running in a test environment
	if isTest {
		// Use environment variable or default to test database
		testDBSource := os.Getenv("DB_TEST_SOURCE")
		if testDBSource != "" {
			return testDBSource
		}
		return "postgresql://root:secret@localhost:5432/musli_test?sslmode=disable"
	}

	// Use environment variable or default to development database
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource != "" {
		return dbSource
	}
	return "postgresql://root:secret@localhost:5432/musli?sslmode=disable"
}

// GetDBDriver returns the database driver
func GetDBDriver() string {
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver != "" {
		return dbDriver
	}
	return "postgres"
}

// PostgresDSN creates a PostgreSQL data source name
func PostgresDSN(host string, port int, user, password, dbname string, sslmode string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}
