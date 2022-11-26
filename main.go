package main

import (
	"context"
	social_graph "github.com/FTN-TwitterClone/grpc-stubs/proto/social_graph"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"social-graph/controller"
	"social-graph/controller/jwt"
	"social-graph/repository/neo4jRepo"
	"social-graph/service"
	"social-graph/tls"
	"social-graph/tracing"
	"time"
)

func main() {
	ctx := context.Background()
	exp, err := tracing.NewExporter()
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := tracing.NewTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	tracer := tp.Tracer("social-graph")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	repositoryNeo4j, err := neo4jRepo.NewRepositoryNeo4j(tracer)
	if err != nil {
		log.Fatal(err)
	}

	socialGraphService := service.NewSocialGraphService(repositoryNeo4j, tracer)

	socialGraphController := controller.NewSocialGraphController(socialGraphService, tracer)
	router := mux.NewRouter()
	router.Use(
		tracing.ExtractTraceInfoMiddleware,
		jwt.ExtractJWTUserMiddleware(tracer),
	)

	router.HandleFunc("/follows", socialGraphController.CreateFollow).Methods("POST")
	router.HandleFunc("/follows", socialGraphController.RemoveFollow).Methods("DELETE")
	router.HandleFunc("/following/{username}", socialGraphController.GetFollowing).Methods("GET")
	router.HandleFunc("/following/{username}/count", socialGraphController.GetNumberOfFollowing).Methods("GET")
	router.HandleFunc("/followers/{username}", socialGraphController.GetFollowers).Methods("GET")
	router.HandleFunc("/followers/{username}/count", socialGraphController.GetNumberOfFollowers).Methods("GET")
	router.HandleFunc("/follows/{from}/{to}", socialGraphController.CheckIfFollowExists).Methods("GET")
	router.HandleFunc("/follows/{from}/{to}", socialGraphController.AcceptRejectFollowRequest).Methods("PATCH")

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	// start server
	srv := &http.Server{
		Addr:      "0.0.0.0:8000",
		Handler:   handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router),
		TLSConfig: tls.GetHTTPServerTLSConfig(),
	}

	go func() {
		log.Println("server starting")

		certFile := os.Getenv("CERT")
		keyFile := os.Getenv("KEY")

		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds := credentials.NewTLS(tls.GetgRPCClientTLSConfig())

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	social_graph.RegisterSocialGraphServiceServer(grpcServer, service.NewgRPCSocialGraphService(tracer, repositoryNeo4j))
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	log.Println("Received terminate, graceful shutdown", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if srv.Shutdown(ctx) != nil {
		log.Fatal("Cannot gracefully shutdown...")
	}
	log.Println("Server stopped")
}
