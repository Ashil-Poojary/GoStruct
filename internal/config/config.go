package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string
	ServerPort     string
	Databases      map[string]DBConfig
	CORS           CORSConfig
	MicrosoftOAuth MicrosoftOAuthConfig
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
	// Load .env first
	if err := godotenv.Load(".env.development"); err != nil {
		log.Printf("Warning: .env.development not loaded: %v", err)
	}

	viper.SetConfigName("config.development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config") // recommended subdirectory
	viper.AddConfigPath(".")      // fallback if not in /config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv() // Allow ENV override

	// Load multiple DB configs into a map
	dbConfigs := make(map[string]DBConfig)
	if err := viper.UnmarshalKey("databases", &dbConfigs); err != nil {
		log.Printf("Warning: failed to parse databases config: %v", err)
	}

	// Inside LoadConfig, before returning Cfg
	microsoftConfig := MicrosoftOAuthConfig{
		ClientID:     viper.GetString("microsoft.client_id"),
		ClientSecret: viper.GetString("microsoft.client_secret"),
		Tenant:       viper.GetString("microsoft.tenant"),
		RedirectURI:  viper.GetString("microsoft.redirect_uri"),
	}

	Cfg = &Config{
		Env:        viper.GetString("env"),
		ServerPort: viper.GetString("PORT"),
		Databases:  dbConfigs,
		CORS: CORSConfig{
			AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
			AllowedMethods:   viper.GetStringSlice("cors.allowed_methods"),
			AllowedHeaders:   viper.GetStringSlice("cors.allowed_headers"),
			AllowCredentials: viper.GetBool("cors.allow_credentials"),
		},
		MicrosoftOAuth: microsoftConfig,
	}

	if Cfg.ServerPort == "" {
		Cfg.ServerPort = "8080" // fallback default
	}

	log.Println("Environment:", Cfg.Env)
	return Cfg, nil
}
