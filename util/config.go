package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	// Database Configuration
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	// Server Configuration
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`

	// AWS Cognito Configuration
	CognitoRegion       string `mapstructure:"COGNITO_REGION"`
	CognitoUserPoolID   string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoClientID     string `mapstructure:"COGNITO_CLIENT_ID"`
	CognitoClientSecret string `mapstructure:"COGNITO_CLIENT_SECRET"`
	CognitoRedirectURL  string `mapstructure:"COGNITO_REDIRECT_URL"`
	CognitoIssuerURL    string `mapstructure:"COGNITO_ISSUER_URL"`

	// // Google OAuth Configuration
	// GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	// GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	// GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`

	// General OAuth Redirect URI
	OAuthRedirectURI string `mapstructure:"OAUTH_REDIRECT_URI"`
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

// GetCognitoClientID returns the Cognito client ID
func GetCognitoClientID() string {
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	if clientID != "" {
		return clientID
	}
	return "6bo0q3c938g1oa0hjggqbdv0b" // Default value
}

// GetCognitoClientSecret returns the Cognito client secret
func GetCognitoClientSecret() string {
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")
	if clientSecret != "" {
		return clientSecret
	}
	return "rhlqr4vp21s2v5rfp7fijltlhe8ha3aj4i5oar561h8hgsvslam" // Default value
}

// GetCognitoRedirectURL returns the Cognito redirect URL
func GetCognitoRedirectURL() string {
	redirectURL := os.Getenv("COGNITO_REDIRECT_URL")
	if redirectURL != "" {
		return redirectURL
	}
	return "http://localhost:8080/callback" // Default value
}

// GetCognitoIssuerURL returns the Cognito issuer URL
func GetCognitoIssuerURL() string {
	issuerURL := os.Getenv("COGNITO_ISSUER_URL")
	if issuerURL != "" {
		return issuerURL
	}
	return "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_177Be0rjJ" // Default value
}

// GetCognitoRegion returns the AWS Cognito region
func GetCognitoRegion() string {
	region := os.Getenv("COGNITO_REGION")
	if region != "" {
		return region
	}
	return "us-east-1" // Default value based on your issuer URL
}

// GetCognitoUserPoolID returns the Cognito user pool ID
func GetCognitoUserPoolID() string {
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	if userPoolID != "" {
		return userPoolID
	}
	return "us-east-1_177Be0rjJ" // Default value extracted from your issuer URL
}

// GetOAuthRedirectURI returns the general OAuth redirect URI
func GetOAuthRedirectURI() string {
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")
	if redirectURI != "" {
		return redirectURI
	}
	// Fallback to Cognito redirect URL if not specifically set
	return GetCognitoRedirectURL()
}
