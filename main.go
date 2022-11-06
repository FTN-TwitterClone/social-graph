package main

import (
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"net/http"
	"os"
	"os/signal"
	"social-graph/controller"
	"social-graph/repository/neo4jRepo"
	"social-graph/service"
	"time"
)

type Neo4jConfiguration struct {
	Url      string
	Username string
	Password string
	Database string
}

func (nc *Neo4jConfiguration) newDriver() neo4j.Driver {
	result, err := neo4j.NewDriver(nc.Url, neo4j.BasicAuth(nc.Username, nc.Password, ""))
	if err != nil {
		panic(err)
	}
	return result
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	configuration := parseConfiguration()
	repositoryNeo4j := neo4jRepo.RepositoryNeo4j{Driver: configuration.newDriver(), Database: configuration.Database}

	socialGraphService := service.NewSocialGraphService(&repositoryNeo4j)

	socialGraphController := controller.NewSocialGraphController(socialGraphService)
	router := mux.NewRouter()

	router.HandleFunc("/api/social-graph/follows", socialGraphController.CreateFollow).Methods("POST")
	router.HandleFunc("/api/social-graph/follows", socialGraphController.RemoveFollow).Methods("DELETE")
	router.HandleFunc("/api/social-graph/following/{username}", socialGraphController.GetFollowing).Methods("GET")
	router.HandleFunc("/api/social-graph/followers/{username}", socialGraphController.GetFollowers).Methods("GET")

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	log.Println("Server listening on port", port)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	log.Println("Received terminate, graceful shutdown", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if server.Shutdown(ctx) != nil {
		log.Fatal("Cannot gracefully shutdown...")
	}
	log.Println("Server stopped")
}

func parseConfiguration() *Neo4jConfiguration {
	return &Neo4jConfiguration{
		Url:      lookupEnvOrGetDefault("NEO4J_URI", "neo4j://localhost:7687"),
		Username: lookupEnvOrGetDefault("NEO4J_USER", "neo4j"),
		Password: lookupEnvOrGetDefault("NEO4J_PASSWORD", "matematika1000"),
		Database: "neo4j",
	}

}

func lookupEnvOrGetDefault(key string, defaultValue string) string {
	if env, found := os.LookupEnv(key); !found {
		return defaultValue
	} else {
		return env
	}
}
