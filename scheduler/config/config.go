package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var config *Config
var m sync.Mutex

type Config struct {
	Env        string     `yaml:"env"`
	App        App        `yaml:"app"`
	HttpServer HttpServer `yaml:"httpServer"`
	Log        Log        `yaml:"log"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Sqlite     Sqlite     `yaml:"sqlite"`
	PostResult PostResult `yaml:"postResult"`
}

type HttpServer struct {
	Port int `yaml:"port"`
}

type Log struct {
	Level string `yaml:"level"`
}

type App struct {
	Name string `yaml:"name"`
}

type Sqlite struct {
	Name            string `yaml:"name"`
	Path            string `yaml:"path"`
	MaxConnections  int    `yaml:"maxConnections"`
	MaxConnIdleTime int    `yaml:"maxConnIdleTime"`
}

type Scheduler struct {
	Timezone string `yaml:"timezone"`
}

type PostResult struct {
	Url    string `yaml:"url"`
	Port   int    `yaml:"port"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
