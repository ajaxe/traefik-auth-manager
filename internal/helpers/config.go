package helpers

import (
	"fmt"
	"log"
	"net/url"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Server struct {
		Port string  `mapstructure:"port"`
	} `mapstructure:"server"`
	OAuth struct {
		Authority    string `mapstructure:"authority"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		CallbackPath string `mapstructure:"callback_path"`
		AppHostURL   string `mapstructure:"app_host_url"`
	} `mapstructure:"oauth"`
	Session struct {
		SessionKey string `mapstructure:"session_key"`
	} `mapstructure:"session"`
	Database struct {
		ConnectionURI string `mapstructure:"connection_uri"`
	} `mapstructure:"database"`
}

func (a AppConfig) OAuthRedirectURL() string {
	p, err := url.JoinPath(a.OAuth.AppHostURL, a.OAuth.CallbackPath)
	if err != nil {
		log.Fatalf("invalid oauth config: %v", err)
	}
	return p
}

func MustLoadDefaultAppConfig() (config AppConfig) {
	config, err := LoadAppConfig(".", "config")
	if err != nil {
		log.Fatalf("failed to get app config: %v", err)
	}
	return
}

func LoadAppConfig(path, name string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetEnvPrefix("app_")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	return
}
