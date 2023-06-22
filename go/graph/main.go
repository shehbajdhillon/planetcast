package graph

import (
	"context"
	"fmt"
	"log"
	"planetcastdev/auth"
	"planetcastdev/database"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

func GenerateServer(queries *database.Queries) *handler.Server {

	gqlConfig := Config{Resolvers: &Resolver{DB: queries}}

	gqlConfig.Directives.LoggedIn = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		if isLoggedIn(ctx) == false {
			return nil, fmt.Errorf("Access Denied")
		}
		return next(ctx)
	}

	var MB int64 = 1 << 20

	gqlServer := handler.New(NewExecutableSchema(gqlConfig))
	gqlServer.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlServer.AddTransport(transport.Options{})
	gqlServer.AddTransport(transport.GET{})
	gqlServer.AddTransport(transport.POST{})
	gqlServer.AddTransport(transport.MultipartForm{MaxUploadSize: 2024 * MB, MaxMemory: 1024 * MB})

	gqlServer.SetQueryCache(lru.New(1000))

	gqlServer.Use(extension.Introspection{})
	gqlServer.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	gqlServer.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		user := auth.FromContext(ctx)
		oc := graphql.GetOperationContext(ctx)
		if user == nil {
			log.Println("Incoming Request:", oc.OperationName, "User: nil")
		} else {
			emailAddr, _ := auth.EmailFromContext(ctx)
			fullName, _ := auth.FullnameFromContext(ctx)

			user, err := queries.GetUserByEmail(ctx, emailAddr)
			if err != nil {
				user, _ = queries.AddUser(ctx, database.AddUserParams{Email: emailAddr, FullName: fullName})
			}
			log.Println("Incoming Request:", oc.OperationName, "User:", user.Email)
		}
		return next(ctx)
	})

	return gqlServer
}

func isLoggedIn(ctx context.Context) bool {
	user := auth.FromContext(ctx)
	return user != nil
}
