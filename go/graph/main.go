package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/email"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/paymentsmiddleware"
	"planetcastdev/storage"
	"planetcastdev/youtubemiddleware"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"go.uber.org/zap"
)

type GraphConnectProps struct {
	Queries  *database.Queries
	Storage  *storage.Storage
	Dubbing  *dubbing.Dubbing
	Logger   *zap.Logger
	Email    *email.Email
	Youtube  *youtubemiddleware.Youtube
	Ffmpeg   *ffmpegmiddleware.Ffmpeg
	Payments *paymentsmiddleware.Payments
}

func Connect(args GraphConnectProps) *handler.Server {

	gqlConfig := Config{Resolvers: &Resolver{
		DB:       args.Queries,
		Storage:  args.Storage,
		Dubbing:  args.Dubbing,
		Logger:   args.Logger,
		Email:    args.Email,
		Youtube:  args.Youtube,
		Ffmpeg:   args.Ffmpeg,
		Payments: args.Payments,
	}}

	logger := args.Logger

	gqlConfig.Directives.LoggedIn = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		if isLoggedIn(ctx) == false {
			return nil, fmt.Errorf("Access Denied")
		}
		return next(ctx)
	}

	gqlConfig.Directives.MemberTeam = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		teamSlugField := obj.(map[string]interface{})["teamSlug"]
		if teamSlugField == nil {
			return nil, fmt.Errorf("Access Denied")
		}
		teamSlug := teamSlugField.(string)
		if memberTeam(ctx, teamSlug, args.Queries) == false {
			return nil, fmt.Errorf("Access Denied")
		}
		return next(ctx)
	}

	gqlConfig.Directives.OwnsProject = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		projectIdField := obj.(map[string]interface{})["projectId"]
		if projectIdField == nil {
			return nil, fmt.Errorf("Access Denied")
		}
		projectId, err := projectIdField.(json.Number).Int64()
		if err != nil || ownsProject(ctx, projectId, args.Queries) == false {
			return nil, fmt.Errorf("Access Denied")
		}
		return next(ctx)
	}

	gqlConfig.Directives.OwnsTransformation = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		transformationIdField := obj.(map[string]interface{})["transformationId"]
		if transformationIdField == nil {
			return nil, fmt.Errorf("Access Denied")
		}
		transformationId, err := transformationIdField.(json.Number).Int64()
		if err != nil || ownsTransformation(ctx, transformationId, args.Queries) == false {
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
	gqlServer.AddTransport(transport.MultipartForm{MaxUploadSize: 1024 * MB, MaxMemory: 1024 * MB})

	gqlServer.SetQueryCache(lru.New(1000))

	gqlServer.Use(extension.Introspection{})
	gqlServer.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	gqlServer.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		user := auth.FromContext(ctx)
		oc := graphql.GetOperationContext(ctx)
		if user == nil {
			logger.Info("Incoming Request", zap.String("operation_name", oc.OperationName), zap.String("user", "nil"))
		} else {
			emailAddr, _ := auth.EmailFromContext(ctx)
			fullName, _ := auth.FullnameFromContext(ctx)

			user, err := args.Queries.GetUserByEmail(ctx, emailAddr)
			if err != nil {
				user, _ = args.Queries.AddUser(ctx, database.AddUserParams{Email: emailAddr, FullName: strings.Title(strings.ToLower(fullName))})
			}
			logger.Info("Incoming Request", zap.String("operation_name", oc.OperationName), zap.String("user", user.Email))
		}
		return next(ctx)
	})

	return gqlServer
}

func isLoggedIn(ctx context.Context) bool {
	user := auth.FromContext(ctx)
	return user != nil
}

func isSuperAdmin(ctx context.Context) bool {
	userEmail, _ := auth.EmailFromContext(ctx)

	if strings.Split(userEmail, "@")[1] == "planetcast.ai" {
		return true
	}

	if userEmail == "shehbaj.dhillon@gmail.com" {
		return true
	}

	return false
}

func memberTeam(ctx context.Context, teamSlug string, queries *database.Queries) bool {
	if !isLoggedIn(ctx) {
		return false
	}

	if isSuperAdmin(ctx) {
		return true
	}

	userEmail, _ := auth.EmailFromContext(ctx)
	user, _ := queries.GetUserByEmail(ctx, userEmail)
	team, _ := queries.GetTeamBySlug(ctx, teamSlug)

	_, err := queries.GetTeamMembershipByTeamIdUserId(ctx, database.GetTeamMembershipByTeamIdUserIdParams{
		UserID: user.ID,
		TeamID: team.ID,
	})

	return err == nil
}

func ownsProject(ctx context.Context, projectId int64, queries *database.Queries) bool {
	project, err := queries.GetProjectById(ctx, projectId)
	if err != nil {
		return false
	}
	team, err := queries.GetTeamById(ctx, project.TeamID)
	return err == nil && memberTeam(ctx, team.Slug, queries)
}

func ownsTransformation(ctx context.Context, transformationId int64, queries *database.Queries) bool {
	transformation, err := queries.GetTransformationById(ctx, transformationId)
	if err != nil {
		return false
	}
	return ownsProject(ctx, transformation.ProjectID, queries)
}
