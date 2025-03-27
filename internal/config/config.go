package config

import (
	"log"
	"path"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"APP" json:"app" yaml:"app"`
	DB    DBConfig    `mapstructure:"DB" json:"db" yaml:"db"`
	JWT   JWTConfig   `mapstructure:"JWT" json:"jwt" yaml:"jwt"`
	OAuth OAuthConfig `mapstructure:"OAUTH" json:"oauth" yaml:"oauth"`
}

type AppConfig struct {
	Addr         string   `mapstructure:"ADDR" json:"addr" yaml:"addr"`
	Domain       string   `mapstructure:"DOMAIN" json:"domain" yaml:"domain"`
	AllowOrigins []string `mapstructure:"ALLOW_ORIGINS" json:"allow_origins" yaml:"allow_origins"`
}

type DBConfig struct {
	File string `mapstructure:"FILE" json:"file" yaml:"file"`
}

type JWTConfig struct {
	Room JWTRoom `mapstructure:"ROOM" json:"room" yaml:"room"`
	User JWTUser `mapstructure:"AUTH" json:"auth" yaml:"auth"`
}

type JWTRoom struct {
	SecretKey string        `mapstructure:"SECRET_KEY" json:"secret_key" yaml:"secret_key"`
	TokenTTL  time.Duration `mapstructure:"TOKEN_TTL" json:"token_ttl" yaml:"token_ttl"`
}

type JWTUser struct {
	SecretKey       string        `mapstructure:"SECRET_KEY" json:"secret_key" yaml:"secret_key"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL" json:"access_token_ttl" yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL" json:"refresh_token_ttl" yaml:"refresh_token_ttl"`
}

type OAuthConfig struct {
	Google OAuthProviderConfig `mapstructure:"GOOGLE" json:"google" yaml:"google"`
	GitHub OAuthProviderConfig `mapstructure:"GITHUB" json:"github" yaml:"github"`
	Yandex OAuthProviderConfig `mapstructure:"YANDEX" json:"yandex" yaml:"yandex"`
}

type OAuthProviderConfig struct {
	ClientID     string `mapstructure:"CLIENT_ID" json:"client_id" yaml:"client_id"`
	ClientSecret string `mapstructure:"CLIENT_SECRET" json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `mapstructure:"REDIRECT_URL" json:"redirect_url" yaml:"redirect_url"`
	UserEndpoint string `mapstructure:"USER_ENDPOINT" json:"user_endpoint" yaml:"user_endpoint"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	viper.SetConfigName(path.Base(file))
	viper.SetConfigType(path.Ext(file)[1:]) // remove dot
	viper.AddConfigPath(path.Dir(file))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("red config file: %v", err)
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("decode config into struct: %v", err)
		return config, err
	}

	return config, nil
}
