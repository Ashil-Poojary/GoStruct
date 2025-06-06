package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string
	ServerPort     string
	DefaultDB      DBConfig `mapstructure:"default_db"`
	ReplicaDB      DBConfig `mapstructure:"replica_db"`
	ProdDB         DBConfig `mapstructure:"prod_db"`
	CORS           CORSConfig
	MicrosoftOAuth MicrosoftOAuthConfig `mapstructure:"microsoft"`
}

type DBConfig struct {
	User      string
	Password  string
	Host      string
	Port      int
	Name      string
	Migration bool
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type MicrosoftOAuthConfig struct {
	ClientID     string
	ClientSecret string
	Tenant       string
	RedirectURI  string
}

var Cfg *Config

func LoadConfig() (*Config, error) {
	// Step 1: Set up Viper to read config file first (without loading .env yet)
	viper.SetConfigName("config.development") // default fallback if no ENV yet
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read config file (ignore error for now to get env)
	_ = viper.ReadInConfig()

	// Step 2: Read 'env' from config or fallback to 'development'
	env := viper.GetString("env")
	if env == "" {
		env = "development"
	}
	log.Printf("[INFO] Environment detected from config: %s", env)

	// Step 3: Load .env file based on detected env
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("[WARNING] Could not load %s file: %v", envFile, err)
	} else {
		log.Printf("[INFO] Loaded environment variables from %s", envFile)
	}

	// Step 4: Now set config name based on detected env and re-read config for overrides
	viper.SetConfigName("config." + env)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file for env '%s': %w", env, err)
	}
	log.Printf("[INFO] Loaded configuration from %s", viper.ConfigFileUsed())

	// Step 5: Bind env vars to override config keys
	bindings := map[string]string{
		"env":                     "ENV",
		"server_port":             "PORT",
		"microsoft.client_id":     "MICROSOFT_CLIENT_ID",
		"microsoft.client_secret": "MICROSOFT_CLIENT_SECRET",
		"microsoft.tenant":        "MICROSOFT_TENANT",
		"microsoft.redirect_uri":  "MICROSOFT_REDIRECT_URI",
		"default_db.user":         "DEFAULT_DB_USER",
		"default_db.password":     "DEFAULT_DB_PASS",
		"replica_db.user":         "REPLICA_DB_USER",
		"replica_db.password":     "REPLICA_DB_PASS",
		"prod_db.user":            "PROD_DB_USER",
		"prod_db.password":        "PROD_DB_PASS",
	}

	for key, envVar := range bindings {
		if err := viper.BindEnv(key, envVar); err != nil {
			log.Printf("[ERROR] Failed to bind env var %s: %v", envVar, err)
		}
	}

	// Step 6: Unmarshal config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Step 7: Fallback default port
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
		log.Println("[INFO] Using default port 8080 as PORT was not set")
	}

	Cfg = &cfg
	log.Printf("[INFO] Final configuration loaded for environment: %s", cfg.Env)
	return Cfg, nil
}
