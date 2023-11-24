package main

import (
	"log"
	"net/http"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/elevenlabsmiddleware"
	"planetcastdev/email"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/graph"
	"planetcastdev/logmiddleware"
	"planetcastdev/openaimiddleware"
	"planetcastdev/paymentsmiddleware"
	"planetcastdev/replicatemiddleware"
	"planetcastdev/storage"
	"planetcastdev/youtubemiddleware"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	production := os.Getenv("PRODUCTION") != ""
	Logger := logmiddleware.Connect(production)

	err := godotenv.Load()
	if err != nil {
		Logger.Warn(".env: Could not find .env file", zap.Error(err))
	} else {
		Logger.Info(".env: Loaded environment variables")
	}

	Payments := paymentsmiddleware.Connect()
	Replicate := replicatemiddleware.Connect(replicatemiddleware.ReplicateConnectProps{Logger: Logger})
	ElevenLabs := elevenlabsmiddleware.Connect(elevenlabsmiddleware.ElevenLabsConnectProps{Logger: Logger})
	OpenAI := openaimiddleware.Connect(openaimiddleware.OpenAIConnectProps{Logger: Logger})
	Email := email.Connect(email.EmailConnectProps{Logger: Logger})
	Ffmpeg := ffmpegmiddleware.Connect(ffmpegmiddleware.FfmpegConnectProps{Logger: Logger})
	Youtube := youtubemiddleware.Connect(youtubemiddleware.YoutubeConnectProps{Logger: Logger, Ffmpeg: Ffmpeg})
	Storage := storage.Connect(storage.StorageConnectProps{Logger: Logger})
	Database := database.Connect(database.DatabaseConnectProps{Logger: Logger})

	Dubbing := dubbing.Connect(
		dubbing.DubbingConnectProps{
			Storage:    Storage,
			Database:   Database,
			Logger:     Logger,
			Ffmpeg:     Ffmpeg,
			Email:      Email,
			Openai:     OpenAI,
			Replicate:  Replicate,
			ElevenLabs: ElevenLabs,
		})

	GqlServer := graph.Connect(graph.GraphConnectProps{
		Dubbing:  Dubbing,
		Storage:  Storage,
		Queries:  Database,
		Logger:   Logger,
		Email:    Email,
		Youtube:  Youtube,
		Ffmpeg:   Ffmpeg,
		Payments: Payments,
	})

	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080", "https://www.planetcast.ai", "https://planetcast.ai", "https://api.planetcast.ai"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            false,
	}).Handler)

	router.Handle("/", GqlServer)

	if production == false {
		Logger.Info("Connect to http://localhost:" + port + " for GraphQL server")
		Logger.Info("connect to http://localhost:" + port + "/playground for GraphQL playground")
		router.Handle("/playground", playground.Handler("GraphQL playground", "/"))
	} else {
		Logger.Info("Connect to https://api.planetcast.ai for GraphQL server")
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
