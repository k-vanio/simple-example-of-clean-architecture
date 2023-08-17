package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/k-vanio/simple-example-of-clean-architecture/configs"
	_ "github.com/k-vanio/simple-example-of-clean-architecture/docs"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/event/handler"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/graph"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/grpc/pb"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/grpc/service"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/web/webserver"
	"github.com/k-vanio/simple-example-of-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pressly/goose/v3"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

// @title           Orders API
// @version         0.1

// @contact.name   Vanio
// @contact.url    https://www.linkedin.com/in/vanio-almeida/
// @contact.email  almeida.vanio@pm.me

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	incorporationPath := filepath.Join("..", "..", "internal", "infra", "database")
	baseFS := os.DirFS(incorporationPath)
	goose.SetBaseFS(baseFS)
	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	orderUseCase := NewOrderUseCase(db, eventDispatcher)

	webServer := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webServer.AddHandler(&webserver.Route{Path: "/orders", Method: http.MethodPost, Action: webOrderHandler.Create})
	webServer.AddHandler(&webserver.Route{Path: "/orders", Method: http.MethodGet, Action: webOrderHandler.List})

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webServer.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*orderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			OrderUseCase: *orderUseCase,
		},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	go http.ListenAndServe(":"+configs.GraphQLServerPort, nil)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	fmt.Println("Pressione Ctrl+C para encerrar o programa.")
	<-signals

	os.Exit(0)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
