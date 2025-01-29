package main

import (
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"library-api-book/internal/config"
	"library-api-book/internal/factory"
	"library-api-book/internal/grpc/client"
	"library-api-book/internal/grpc/handlers"
	"library-api-book/internal/routes"
	"library-api-book/pkg/database"
	"library-api-book/proto/book"
)

func main() {
	config.LoadConfig()

	psqlDB, err := database.NewPqSQLClient()
	if err != nil {
		log.Fatal("Could not connect to PqSQL:", err)
	}

	provider := factory.InitFactory(psqlDB)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runGRPCServer(provider)
	}()

	go func() {
		defer wg.Done()
		runHTTPServer(provider)
	}()

	wg.Wait()
}

func runGRPCServer(provider *factory.Provider) {
	listener, err := net.Listen("tcp", ":"+config.ENV.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", config.ENV.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()

	bookHandler := handlers.NewBookHandler(provider.BookService)
	book.RegisterBookServiceServer(grpcServer, bookHandler)

	log.Printf("gRPC server running on port %s\n", config.ENV.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func runHTTPServer(provider *factory.Provider) {
	authClient, err := client.NewAuthClient(config.ENV.UserGRPC)
	if err != nil {
		log.Fatalf("Failed to initialize auth client: %v", err)
	}
	defer authClient.Close()

	router := routes.RegisterRoutes(provider, authClient)
	log.Printf("REST API server running on port %s\n", config.ENV.ServerPort)
	log.Fatal(router.Run(":" + config.ENV.ServerPort))
}
