package main

import (
	"github.com/spf13/viper"
	core2 "github.com/zytell3301/tg-users-service/internal/core"
	"github.com/zytell3301/tg-users-service/internal/handlers/grpcHandlers"
	"github.com/zytell3301/tg-users-service/internal/repository"
	"github.com/zytell3301/tg-users-service/pkg/UsersService"
	uuid_generator "github.com/zytell3301/uuid-generator"
	"google.golang.org/grpc"
	"log"
	"net"
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
	configs := configs{}
	configs.repositoryConfigs = loadRepositoryConfigs()
	configs.serviceConfigs = loadServiceConfigs()
	uuidGenerator, err := uuid_generator.NewGenerator(configs.serviceConfigs.uuidSpace)
	switch err != nil {
	case true:
		log.Fatalf("An error occurred while initiating uuid generator instance. Error message: %v", err)
	}
	repo, err := repository.NewUsersRepository(configs.repositoryConfigs.hosts, configs.repositoryConfigs.keyspace, uuidGenerator)
	usersCore := core2.NewUsersCore(repo)
	grpcHandler := grpcHandlers.NewHandler(usersCore)
	listener, err := net.Listen("tcp", configs.serviceConfigs.nodeIp+":"+configs.serviceConfigs.servicePort)
	grpcServer := grpc.NewServer()
	UsersService.RegisterUsersServiceServer(grpcServer, grpcHandler)
	err = grpcServer.Serve(listener)
	switch err != nil {
	case true:
		log.Fatalf("An error occurred while starting grpc servcice. Error message: %v", err)
	}
}

func loadConfig(config string) *viper.Viper {
	cfg := viper.New()
	cfg.AddConfigPath(ProjectRoot + "/configs")
	cfg.SetConfigName(config)
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

func loadServiceConfigs() (config serviceConfigs) {
	cfg := loadConfig("service")
	config.nodeIp = cfg.GetString("node-ip")
	config.servicePort = cfg.GetString("service-port")
	config.uuidSpace = cfg.GetString("uuid-space")
	return
}
