package config

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Application holds application configurations
type Application struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DebugPort       int    `yaml:"debug_port"`
	GracefulTimeout int    `yaml:"graceful_timeout"`
}

var appOnce = sync.Once{}
var appConfig *Application

// loadApp loads config from path
func loadApp(fileName string) error {
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}
	viper.AutomaticEnv()

	appConfig = &Application{
		Host:            viper.GetString("app.host"),
		GracefulTimeout: viper.GetInt("app.graceful_timeout"),
		Port:            viper.GetInt("app.port"),
		DebugPort:       viper.GetInt("app.debug_port"),
	}

	log.Println("app config ", appConfig)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if appConfig != nil {
			appConfig.Host = viper.GetString("app.host")
			appConfig.Port = viper.GetInt("app.port")
			appConfig.GracefulTimeout = viper.GetInt("app.graceful_timeout")
			appConfig.DebugPort = viper.GetInt("app.debug_port")
		}
		log.Printf("Config file changed, event: %v", e.Name)
	})

	return nil
}

// GetApp returns application config
func GetApp(fileName string) *Application {
	appOnce.Do(func() {
		loadApp(fileName)
	})

	return appConfig
}
