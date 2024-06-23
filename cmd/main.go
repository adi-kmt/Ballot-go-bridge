package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/adi-kmt/ai-streak-backend-go/internal/controllers"
	"github.com/adi-kmt/ai-streak-backend-go/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	proto.RegisterAuthServiceServer(s, &controllers.Server{})
	// Serve gRPC server
	log.Println("Serving gRPC on connection ")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	// Register Greeter
	err = proto.RegisterAuthServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
