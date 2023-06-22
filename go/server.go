package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"planetcastdev/database"
	"planetcastdev/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	production := os.Getenv("PRODUCTION")

	err := godotenv.Load()
	if err != nil {
		log.Println(".env: Could not find .env file", err.Error())
	} else {
		log.Println(".env: Loaded environment variables")
	}

	Database := database.Connect()
	fmt.Print(Database)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	if production == "" {
		http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
		log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	}
	http.Handle("/", srv)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
