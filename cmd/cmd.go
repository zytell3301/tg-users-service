package main

import (
	"github.com/spf13/viper"
	"log"
)

const ProjectRoot = "./.."

type configs struct {
	repositoryConfigs repositoryConfigs
	serviceConfigs    serviceConfigs
}

type serviceConfigs struct {
	nodeIp      string
	servicePort string
	uuidSpace   string
}

type repositoryConfigs struct {
	hosts    []string
	keyspace string
}

func main() {

}

func loadConfig(config string) *viper.Viper {
	cfg := viper.New()
	cfg.AddConfigPath(ProjectRoot + "/configs")
	cfg.SetConfigName("configs")
	cfg.SetConfigType("yaml")
	err := cfg.ReadInConfig()
	switch err != nil {
	case true:
		log.Fatalf("An error occurred while reading configs. Error message: %v", err)
	}
	return cfg
}

func loadRepositoryConfigs() (config repositoryConfigs) {
	cfg := loadConfig("repository")
	config.hosts = cfg.GetStringSlice("hosts")
	config.keyspace = cfg.GetString("keyspace")
	return
}
