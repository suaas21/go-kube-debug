package app

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Application holds application configurations
type Application struct {
	OCAgentHost string

	WardenPort      int
	WardenDebugPort int
	MemberSVC       string

	MemberPort int
	SettingSVC string

	SettingPort   int
	SettlementSVC string

	SettlementPort int
	CoreSVC        string

	CorePort int

	Service string
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
		OCAgentHost: viper.GetString("tracing.oc_agent_host"),

		WardenPort:      viper.GetInt("warden.port"),
		WardenDebugPort: viper.GetInt("warden.debug_port"),
		MemberSVC:       viper.GetString("warden.member_svc"),

		MemberPort: viper.GetInt("member.port"),
		SettingSVC: viper.GetString("member.setting_svc"),

		SettingPort:   viper.GetInt("setting.port"),
		SettlementSVC: viper.GetString("setting.settlement_svc"),

		SettlementPort: viper.GetInt("settlement.port"),
		CoreSVC:        viper.GetString("settlement.core_svc"),

		CorePort: viper.GetInt("core.port"),
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
