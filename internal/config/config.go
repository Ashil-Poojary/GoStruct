package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string               `mapstructure:"env"`
	ServerPort     string               `mapstructure:"server_port"`
	DefaultDB      DBConfig             `mapstructure:"default_db"`
	ReplicaDB      DBConfig             `mapstructure:"replica_db"`
	ProdDB         DBConfig             `mapstructure:"prod_db"`
	CORS           CORSConfig           `mapstructure:"cors"`
	MicrosoftOAuth MicrosoftOAuthConfig `mapstructure:"microsoft"`
}

type DBConfig struct {
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Name      string `mapstructure:"name"`
	Migration bool   `mapstructure:"migration"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type MicrosoftOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Tenant       string `mapstructure:"tenant"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}

var Cfg *Config

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load base config
	viper.SetConfigName("config.development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Determine environment (default to development)
	env := viper.GetString("env")
	if env == "" {
		env = "development"
	}

	// Load .env.{env}
	envFile := ".env." + env
	if err := godotenv.Overload(envFile); err != nil {
		log.Printf("Warning: Could not load %s: %v", envFile, err)
	}

	// Load config.{env}.yaml
	viper.SetConfigName("config." + env)
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read %s config: %w", env, err)
	}

	// Bind common environment variables
	envBindings := map[string]string{
		"server_port":             "PORT",
		"default_db.user":         "DEFAULT_DB_USER",
		"default_db.password":     "DEFAULT_DB_PASS",
		"default_db.name":         "DEFAULT_DB_NAME",
		"replica_db.user":         "REPLICA_DB_USER",
		"replica_db.password":     "REPLICA_DB_PASS",
		"replica_db.name":         "REPLICA_DB_NAME",
		"prod_db.user":            "PROD_DB_USER",
		"prod_db.password":        "PROD_DB_PASS",
		"prod_db.name":            "PROD_DB_NAME",
		"microsoft.client_id":     "MICROSOFT_CLIENT_ID",
		"microsoft.client_secret": "MICROSOFT_CLIENT_SECRET",
		"microsoft.tenant":        "MICROSOFT_TENANT",
		"microsoft.redirect_uri":  "MICROSOFT_REDIRECT_URI",
	}

	for key, env := range envBindings {
		if err := viper.BindEnv(key, env); err != nil {
			log.Printf("Warning: failed to bind env var %s: %v", env, err)
		}
	}

	// Unmarshal final config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Fallbacks and validation
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	Cfg = &cfg
	return Cfg, nil
}
