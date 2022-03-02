package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
	ErrorReporter "github.com/zytell3301/tg-error-reporter"
	core2 "github.com/zytell3301/tg-users-service/internal/core"
	"github.com/zytell3301/tg-users-service/internal/errorReporter"
	"github.com/zytell3301/tg-users-service/internal/handlers/grpcHandlers"
	"github.com/zytell3301/tg-users-service/internal/repository"
	"github.com/zytell3301/tg-users-service/pkg/CertGen"
	"github.com/zytell3301/tg-users-service/pkg/UsersService"
	uuid_generator "github.com/zytell3301/uuid-generator"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

const ProjectRoot = "."

type configs struct {
	repositoryConfigs repositoryConfigs
	serviceConfigs    serviceConfigs
}

type serviceConfigs struct {
	nodeIp      string
	servicePort string
	uuidSpace   string
	serviceId   string
	instanceId  string
}

type repositoryConfigs struct {
	hosts             []string
	keyspace          string
	port              int
	consistencyLevels repository.ConsistencyLevels
}

func main() {
	configs := configs{}
	configs.repositoryConfigs = loadRepositoryConfigs()
	configs.serviceConfigs = loadServiceConfigs()
	errorReporter.InitiateReporter(configs.serviceConfigs.instanceId, configs.serviceConfigs.serviceId, ErrorReporter.DefaultReporter{})
	uuidGenerator := newUuidGenerator(configs.serviceConfigs.uuidSpace)
	repo := newUsersRepo(configs.repositoryConfigs.hosts, configs.repositoryConfigs.keyspace, uuidGenerator, configs.repositoryConfigs.consistencyLevels)
	certGen := newCertgen()
	usersCore := core2.NewUsersCore(repo, certGen)
	grpcHandler := grpcHandlers.NewHandler(usersCore)
	listener := newListener(configs)
	grpcServer := grpc.NewServer()
	UsersService.RegisterUsersServiceServer(grpcServer, grpcHandler)
	fmt.Println("Serving grpc server")
	err := grpcServer.Serve(listener)
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

func newUuidGenerator(space string) *uuid_generator.Generator {
	fmt.Println("Creating uuid generator instance...")
	uuidGenerator, err := uuid_generator.NewGenerator(space)
	switch err != nil {
	case true:
		log.Fatalf("An error occurred while initiating uuid generator instance. Error message: %v", err)
	}
	fmt.Println("Uuid generator instance created successfully")
	return uuidGenerator
}

func newUsersRepo(hosts []string, keyspace string, uuidGenerator *uuid_generator.Generator, consistencyLevels repository.ConsistencyLevels) repository.Repository {
	fmt.Println("Creating new users repository instance...")
	repo, err := repository.NewUsersRepository(hosts, keyspace, uuidGenerator, consistencyLevels)
	switch err != nil {
	case true:
		log.Fatalf("An error occurred while creating users repository. Error message: %v", err)
	}
	fmt.Println("Users repository created successfully")
	return repo
}

func newCertgen() CertGen.CertGen {
	fmt.Println("Creating certGen instance...")
	certGen, err := CertGen.NewCertGenerator(getCertificate(), getCertificateKey())
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while creating certGen instance. Error message: %s", err.Error()))
	}
	fmt.Println("CertGen instance created successfully")
	return certGen
}

func newListener(configs configs) net.Listener {
	listener, err := net.Listen("tcp", configs.serviceConfigs.nodeIp+":"+configs.serviceConfigs.servicePort)
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while creating tcp listener. Error message: %s", err.Error()))
	}
	return listener
}

func loadRepositoryConfigs() (config repositoryConfigs) {
	fmt.Println("Loading repository configs")
	cfg := loadConfig("repository")
	config.hosts = cfg.GetStringSlice("hosts")
	config.keyspace = cfg.GetString("keyspace")
	consistencyLevels := cfg.GetStringMapString("consistency-levels")
	config.consistencyLevels.NewUser = parseConsistencyLevel(consistencyLevels["new-user"])
	config.consistencyLevels.GetUserByPhone = parseConsistencyLevel(consistencyLevels["get-user-by-phone"])
	config.consistencyLevels.GetSecurityCode = parseConsistencyLevel(consistencyLevels["get-security-code"])
	config.consistencyLevels.GetUserByUsername = parseConsistencyLevel(consistencyLevels["get-user-by-username"])
	config.consistencyLevels.DoesUserExists = parseConsistencyLevel(consistencyLevels["does-user-exists"])
	config.consistencyLevels.RecordSecurityCode = parseConsistencyLevel(consistencyLevels["record-security-code"])
	config.consistencyLevels.DeleteUser = parseConsistencyLevel(consistencyLevels["delete-user"])
	config.consistencyLevels.UpdateUsername = parseConsistencyLevel(consistencyLevels["update-username"])
	config.consistencyLevels.DoesUsernameExists = parseConsistencyLevel(consistencyLevels["does-username-exists"])
	config.port = cfg.GetInt("port")
	fmt.Println("Repository config loaded successfully")
	return
}

func parseConsistencyLevel(level string) gocql.Consistency {
	switch level {
	case "ALL":
		return gocql.All
	case "ONE":
		return gocql.One
	case "TWO":
		return gocql.Two
	case "THREE":
		return gocql.Three
	case "ANY":
		return gocql.Any
	case "QUORUM":
		return gocql.Quorum
	case "EACH-QUORUM":
		return gocql.EachQuorum
	case "LOCAL-ONE":
		return gocql.LocalOne
	case "LOCAL-QUORUM":
		return gocql.LocalQuorum
	default:
		panic(fmt.Sprintf("Defined consistency level is not valid. Expected:ALL,ANY,ONE,TWO,THREE,QUORUM,LOCAL-QUORUM,EACH,QUORUM,LOCAL-ONE, got: %v", level))
	}
}

func loadServiceConfigs() (config serviceConfigs) {
	fmt.Println("Loading service configs")
	cfg := loadConfig("service")
	config.nodeIp = cfg.GetString("node-ip")
	config.servicePort = cfg.GetString("service-port")
	config.uuidSpace = cfg.GetString("uuid-space")
	config.serviceId = cfg.GetString("service-id")
	config.instanceId = cfg.GetString("instance-id")
	fmt.Println("Service configs loaded successfully")
	return
}

func getCertificate() []byte {
	fmt.Println("Service root certificate is being loaded")
	file, err := os.Open("./auth-certificates/certificate.pem")
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while opening root certificate. Error message: %s", err.Error()))
	}
	cert, err := io.ReadAll(file)
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while reading root certificate. Error messge: %s", err.Error()))
	}
	fmt.Println("Service root certificate loaded successfully")
	return cert
}

func getCertificateKey() []byte {
	fmt.Println("Service certificate private key is being loaded")
	file, err := os.Open("./auth-certificates/key.pem")
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while opening certificate key file. Error message: %s", err.Error()))
	}
	cert, err := io.ReadAll(file)
	switch err != nil {
	case true:
		panic(fmt.Sprintf("An error occurred while reading certificate key file. Error message: %s", err.Error()))
	}
	fmt.Println("Service certificate private key loaded successfully")
	return cert
}
