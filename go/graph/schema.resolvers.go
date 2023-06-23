package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"fmt"
	"planetcastdev/auth"
	"planetcastdev/database"

	"github.com/99designs/gqlgen/graphql"
)

// CreateTeam is the resolver for the createTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, slug string, name string, teamType database.TeamType) (database.Team, error) {
	email, _ := auth.EmailFromContext(ctx)
	user, _ := r.DB.GetUserByEmail(ctx, email)

	team, _ := r.DB.CreateTeam(ctx, database.CreateTeamParams{
		Slug:     slug,
		Name:     name,
		TeamType: teamType,
	})

	r.DB.AddTeamMembership(ctx, database.AddTeamMembershipParams{
		TeamID:         team.ID,
		UserID:         user.ID,
		MembershipType: database.MembershipTypeOWNER,
	})

	return team, nil
}

// CreateProject is the resolver for the createProject field.
func (r *mutationResolver) CreateProject(ctx context.Context, teamID int64, title string, sourceLanguage database.SupportedLanguage, targetLanguage database.SupportedLanguage, sourceMedia graphql.Upload) (database.Project, error) {
	panic(fmt.Errorf("not implemented: CreateProject - createProject"))
}

// GetTeams is the resolver for the getTeams field.
func (r *queryResolver) GetTeams(ctx context.Context) ([]database.Team, error) {
	teams := []database.Team{}
	email, _ := auth.EmailFromContext(ctx)
	user, _ := r.DB.GetUserByEmail(ctx, email)
	memberships, _ := r.DB.GetTeamMemebershipsByUserId(ctx, user.ID)
	for _, mem := range memberships {
		team, err := r.DB.GetTeamById(ctx, mem.TeamID)
		if err == nil {
			teams = append(teams, team)
		}
	}
	return teams, nil
}

// Created is the resolver for the created field.
func (r *teamResolver) Created(ctx context.Context, obj *database.Team) (string, error) {
	return obj.Created.String(), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Team returns TeamResolver implementation.
func (r *Resolver) Team() TeamResolver { return &teamResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
