package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"fmt"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/utils"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
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
func (r *mutationResolver) CreateProject(ctx context.Context, teamSlug string, title string, sourceLanguage database.SupportedLanguage, sourceMedia graphql.Upload) (database.Project, error) {
	team, _ := r.DB.GetTeamBySlug(ctx, teamSlug)

	identifier := strings.Split(sourceMedia.Filename, ".mp4")[0] + uuid.NewString()
	fileName := identifier + ".mp4"

	project, _ := r.DB.CreateProject(ctx, database.CreateProjectParams{
		TeamID:         team.ID,
		Title:          title,
		SourceLanguage: sourceLanguage,
		SourceMedia:    fileName,
	})

	r.Storage.Upload(fileName, sourceMedia.File)

	newCtx := context.Background()

	go dubbing.CreateTransformation(newCtx, r.DB, dubbing.CreateTransformationParams{
		ProjectID:      project.ID,
		TargetLanguage: sourceLanguage,
		FileName:       fileName,
		File:           sourceMedia.File,
		IsSource:       true,
	})

	return project, nil
}

// DeleteProject is the resolver for the deleteProject field.
func (r *mutationResolver) DeleteProject(ctx context.Context, projectID int64) (database.Project, error) {
	project, _ := r.DB.DeleteProjectById(ctx, projectID)
	return project, nil
}

// CreateTranslation is the resolver for the createTranslation field.
func (r *mutationResolver) CreateTranslation(ctx context.Context, projectID int64, targetLanguage database.SupportedLanguage) (database.Transformation, error) {
	// fetch source transcript for the project
	sourceTransformation, err := r.DB.GetSourceTransformationByProjectId(ctx, projectID)
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Project Not Processed!")
	}

	// if target transformation already exists, return that
	existingTransformation, err := r.DB.GetTransformationByProjectIdTargetLanguage(ctx, database.GetTransformationByProjectIdTargetLanguageParams{
		ProjectID:      projectID,
		TargetLanguage: targetLanguage,
	})
	if err == nil {
		return existingTransformation, nil
	}

	identifier := fmt.Sprintf("%d-%s-%s", sourceTransformation.ProjectID, utils.GetCurrentDateTimeString(), targetLanguage)
	newFileName := identifier + "_dubbed.mp4"

	// create empty transformation in target language, if target transformation already exists, return that
	newTransformation, _ := r.DB.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      projectID,
		TargetLanguage: targetLanguage,
		TargetMedia:    newFileName,
		Transcript:     pqtype.NullRawMessage{Valid: false, RawMessage: nil},
		IsSource:       false,
	})

	newCtx := context.Background()
	go dubbing.CreateTranslation(newCtx, r.DB, sourceTransformation, newTransformation, identifier)

	return newTransformation, nil
}

// Transformations is the resolver for the transformations field.
func (r *projectResolver) Transformations(ctx context.Context, obj *database.Project, transformationID *int64) ([]database.Transformation, error) {
	transformations := []database.Transformation{}

	if transformationID != nil {
		transformation, _ := r.DB.GetTransformationByTransformationIdProjectId(ctx, database.GetTransformationByTransformationIdProjectIdParams{
			ID:        *transformationID,
			ProjectID: obj.ID,
		})
		transformations = []database.Transformation{transformation}
	} else {
		transformations, _ = r.DB.GetTransformationsByProjectId(ctx, obj.ID)
	}

	filteredTransformation := []database.Transformation{}
	for _, t := range transformations {
		if len(t.TargetMedia) > 0 {
			t.TargetMedia = r.Storage.GetFileLink(t.TargetMedia)
		}
		filteredTransformation = append(filteredTransformation, t)
	}

	return filteredTransformation, nil
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

// GetTeamByID is the resolver for the getTeamById field.
func (r *queryResolver) GetTeamByID(ctx context.Context, teamSlug string) (database.Team, error) {
	team, _ := r.DB.GetTeamBySlug(ctx, teamSlug)
	return team, nil
}

// Created is the resolver for the created field.
func (r *teamResolver) Created(ctx context.Context, obj *database.Team) (string, error) {
	return obj.Created.String(), nil
}

// Projects is the resolver for the projects field.
func (r *teamResolver) Projects(ctx context.Context, obj *database.Team, projectID *int64) ([]database.Project, error) {
	projects := []database.Project{}

	if projectID != nil {
		project, _ := r.DB.GetProjectByProjectIdTeamId(ctx, database.GetProjectByProjectIdTeamIdParams{
			ID:     *projectID,
			TeamID: obj.ID,
		})
		projects = []database.Project{project}
	} else {
		projects, _ = r.DB.GetProjectsByTeamId(ctx, obj.ID)
	}

	filteredProject := []database.Project{}
	for _, p := range projects {
		p.SourceMedia = r.Storage.GetFileLink(p.SourceMedia)
		filteredProject = append(filteredProject, p)
	}

	return filteredProject, nil
}

// Transcript is the resolver for the transcript field.
func (r *transformationResolver) Transcript(ctx context.Context, obj *database.Transformation) (string, error) {
	jsonBytes := obj.Transcript.RawMessage
	return string(jsonBytes), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Project returns ProjectResolver implementation.
func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Team returns TeamResolver implementation.
func (r *Resolver) Team() TeamResolver { return &teamResolver{r} }

// Transformation returns TransformationResolver implementation.
func (r *Resolver) Transformation() TransformationResolver { return &transformationResolver{r} }

type mutationResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
type transformationResolver struct{ *Resolver }
