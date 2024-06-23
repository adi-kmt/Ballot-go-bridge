package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/adi-kmt/ai-streak-backend-go/internal/controllers"
	"github.com/adi-kmt/ai-streak-backend-go/internal/entities"
	"github.com/adi-kmt/ai-streak-backend-go/internal/injection"
	"github.com/adi-kmt/ai-streak-backend-go/internal/jwt"
	"github.com/adi-kmt/ai-streak-backend-go/proto"
	"github.com/gorilla/websocket"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	userService, votingService := injection.Injector()

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	proto.RegisterAuthServiceServer(s, &controllers.Server{
		Service: userService,
	})
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

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// writing a custom websocket handler for the voting
	http.HandleFunc("/v1/vote", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		token := r.Header.Get("BearerToken")
		user, err0 := jwt.ParseTokenAndGetClaims(token)
		if err0 != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err0.Message))
			return
		}

		for {
			_, message, err1 := conn.ReadMessage()
			if err1 != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error in reading the message"))
				log.Println("WebSocket connection closed:", err1)
				return
			}
			//Trying to add vote, if successful close connection
			err0 := votingService.AddVote(user, string(message))
			if err0 == nil {
				response := []byte("Thamk you for voting")
				conn.WriteMessage(websocket.TextMessage, response)
				log.Println("WebSocket connection closed:", err)
				return
			}
		}
	})

	http.HandleFunc("/v1/dashboard", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		token := r.Header.Get("BearerToken")
		_, err0 := jwt.ParseTokenAndGetClaims(token)
		if err0 != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err0.Message))
			return
		}
		subscription := entities.NewLeaderBoardSubscription()
		go votingService.GetCurrentVoteSapshot(subscription)

		for {
			select {
			case leaderboard := <-subscription.UpdateChan:
				conn.WriteJSON(leaderboard)
			}
		}
	})

	gwServer := &http.Server{
		Addr: ":9091",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
				mux.ServeHTTP(resp, req)
			} else if req.URL.Path == "/v1/vote" {
				http.DefaultServeMux.ServeHTTP(resp, req)
			} else if req.URL.Path == "/v1/dashboard" {
				http.DefaultServeMux.ServeHTTP(resp, req)
			} else {
				// Handle other HTTP requests here
				resp.WriteHeader(http.StatusNotFound)
				resp.Write([]byte("Not Found"))
			}
		}),
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
