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
	OCAgentHost     string `yaml:"oc_agent_host"`
	RequestPort     int    `yaml:"request_port"`
	Svc             string `yaml:"svc"`
	ReqSvc          string `yaml:"req_svc"`
	ReqPort         int    `yaml:"req_port"`
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
		OCAgentHost:     viper.GetString("tracing.oc_agent_host"),
		RequestPort:     viper.GetInt("request.port"),
		Svc:             viper.GetString("request.svc"),
		ReqSvc:          viper.GetString("req.svc"),
		ReqPort:         viper.GetInt("req.port"),
	}

	log.Println("app config ", appConfig)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
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
