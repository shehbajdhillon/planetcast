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

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(".env: Could not find .env file", err.Error())
	} else {
		log.Println(".env: Loaded environment variables")
	}

	Database := database.Connect()
	fmt.Print(Database)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
