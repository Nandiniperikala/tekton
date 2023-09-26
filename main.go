package main

import (
	auth "com.nandini.icici/auth"
	"com.nandini.icici/controllers"
	"github.com/asim/go-micro/v3"
	rabbitmq "com.nandini.icici/rabbitmq"
	"github.com/micro/micro/v3/service/logger"
	eureka "com.nandini.icici/eurekaregistry"
	"github.com/google/uuid"
	"com.nandini.icici/handler"
	"com.nandini.icici/migrate"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
	mhttp "github.com/go-micro/plugins/v3/server/http"
   "github.com/gorilla/mux"
	app "com.nandini.icici/config"
)

var configurations eureka.RegistrationVariables

func main() {
	defer cleanup()
	app.Setconfig()
	migrate.MigrateAndCreateDatabase()
	auth.SetClient()
	handler.InitializeDb()
	service_registry_url :=app.GetVal("GO_MICRO_SERVICE_REGISTRY_URL")
	InstanceId := "icici:"+uuid.New().String()
	configurations = eureka.RegistrationVariables {ServiceRegistryURL:service_registry_url,InstanceId:InstanceId}
	port :=app.GetVal("GO_MICRO_SERVICE_PORT")
	srv := micro.NewService(
		micro.Server(mhttp.NewServer()),
    )
	opts1 := []micro.Option{
		micro.Name("icici"),
		micro.Version("latest"),
		micro.Address(":"+port),
	}
	srv.Init(opts1...)
	r := mux.NewRouter().StrictSlash(true)
	registerRoutes(r)		
	var handlers http.Handler = r
	
	go eureka.ManageDiscovery(configurations)
	go rabbitmq.ConsumerAndhraToIcici() 

    if err := micro.RegisterHandler(srv.Server(), handlers); err != nil {
		logger.Fatal(err)
	}
	
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

func cleanup(){
	eureka.Cleanup(configurations)
}

func registerRoutes(router *mux.Router) {
	registerControllerRoutes(controllers.EventController{}, router)
}

func registerControllerRoutes(controller controllers.Controller, router *mux.Router) {
	controller.RegisterRoutes(router)
}
