package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	jaegar "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/vishnusunil243/OrderService/db"
	"github.com/vishnusunil243/OrderService/initializer"
	"github.com/vishnusunil243/OrderService/service"
	servicediscovery "github.com/vishnusunil243/OrderService/service_discovery"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}
	cartConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		cartConn.Close()
	}()
	cartRes := pb.NewCartServiceClient(cartConn)
	service.CartClient = cartRes
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatal("failed to listen on port 8084")
	}
	log.Printf("order server listeningon port 8084")
	go func() {
		time.Sleep(5 * time.Second)
		servicediscovery.RegisterService()
	}()
	healthService := &service.HealthChecker{}
	grpc_health_v1.RegisterHealthServer(server, healthService)
	tracer, closer := InitTracer()
	defer closer.Close()
	service.RetrieveTracer(tracer)
	if err = server.Serve(listener); err != nil {
		log.Fatal("failed to listen on port 8084")
	}
}
func InitTracer() (tracer opentracing.Tracer, closer io.Closer) {
	jaegarEndpoint := "http://localhost:14268/api/tracer"
	cfg := &config.Configuration{
		ServiceName: "order-service",
		Sampler: &config.SamplerConfig{
			Type:  jaegar.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaegarEndpoint,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("updated")
	return
}
