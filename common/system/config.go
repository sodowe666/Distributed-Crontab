package system

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Cfg struct {
	ConfigPath string
}

func SetupConfig(env string, configPath string) error {
	config := new(Cfg)
	config.ConfigPath = configPath
	if err := config.initConfig(env); err != nil {
		return err
	}
	//监听配置文件
	config.watchConfig()
	return nil
}

func (c *Cfg) initConfig(env string) error {

	//判断有没有输入配置文件 ，没有则解析默认配置文件
	if c.ConfigPath != "" {
		viper.SetConfigFile(c.ConfigPath)
	} else {
		viper.AddConfigPath("./conf/" + env)
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Cfg) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("Config file changed: %s", in.Name)
	})
}
