package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"social-graph/controller"
	"social-graph/repository/neo4jRepo"
	"social-graph/service"
	"social-graph/tls"
	"time"
)

func main() {
	repositoryNeo4j, err := neo4jRepo.NewRepositoryNeo4j()
	if err != nil {
		log.Fatal(err)
	}

	socialGraphService := service.NewSocialGraphService(repositoryNeo4j)

	socialGraphController := controller.NewSocialGraphController(socialGraphService)
	router := mux.NewRouter()

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
