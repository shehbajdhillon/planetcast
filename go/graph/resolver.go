package graph

import (
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/email"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/paymentsmiddleware"
	"planetcastdev/storage"
	"planetcastdev/youtubemiddleware"

	"go.uber.org/zap"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *database.Queries
	Storage  *storage.Storage
	Dubbing  *dubbing.Dubbing
	Logger   *zap.Logger
	Email    *email.Email
	Youtube  *youtubemiddleware.Youtube
	Ffmpeg   *ffmpegmiddleware.Ffmpeg
	Payments *paymentsmiddleware.Payments
}
