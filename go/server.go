package main

import (
	"log"
	"net/http"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/graph"
	"planetcastdev/storage"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	err := godotenv.Load()
	if err != nil {
		log.Println(".env: Could not find .env file", err.Error())
	} else {
		log.Println(".env: Loaded environment variables")
	}

	production := os.Getenv("PRODUCTION")

	Storage := storage.Connect()

	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080", "https://www.planetcast.ai", "https://planetcast.ai", "https://api.planetcast.ai"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            false,
	}).Handler)

	Database := database.Connect()

	srv := graph.GenerateServer(Database, Storage)
	router.Handle("/", srv)

	if production == "" {
		log.Println("Connect to http://localhost:" + port + " for GraphQL server")
		router.Handle("/playground", playground.Handler("GraphQL playground", "/"))
		log.Println("connect to http://localhost:" + port + "/playground for GraphQL playground")
	} else {
		log.Println("Connect to https://api.planetcast.ai for GraphQL server")
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
