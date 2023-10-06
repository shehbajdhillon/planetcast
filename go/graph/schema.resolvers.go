package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"fmt"
	"io"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/graph/model"
	"planetcastdev/utils"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
	"go.uber.org/zap"
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
func (r *mutationResolver) CreateProject(ctx context.Context, teamSlug string, title string, sourceMedia *graphql.Upload, youtubeLink *string, uploadOption model.UploadOption, initialTargetLanguage *model.SupportedLanguage, initialLipSync bool) (database.Project, error) {
	team, _ := r.DB.GetTeamBySlug(ctx, teamSlug)

	// check if file upload or youtube
	// if youtube link, validate the link, if valid, start download

	var file io.ReadSeeker
	var identifier string
	var fileName string

	if uploadOption == model.UploadOptionYoutubeLink {
		_, err := r.Youtube.GetVideoInfo(*youtubeLink)
		if err != nil {
			return database.Project{}, fmt.Errorf("Error processing YouTube video: %s", err.Error())
		}
	}

	project, _ := r.DB.CreateProject(ctx, database.CreateProjectParams{
		TeamID:      team.ID,
		Title:       title,
		SourceMedia: "",
	})

	user := auth.FromContext(ctx)
	newCtx := context.Background()
	newCtx = auth.AttachContext(newCtx, user)

	go func(context context.Context) {

		randomString := uuid.NewString()

		if uploadOption == model.UploadOptionYoutubeLink {
			youtubeFile, youtubeFileName, err := r.Youtube.Download(*youtubeLink)

			if err != nil {
				r.Logger.Info("Could not download youtube video for project", zap.Error(err), zap.Int64("project_id", project.ID), zap.String("youtube_url", *youtubeLink))
				return
			}

			file = youtubeFile
			fileName = strings.ReplaceAll(youtubeFileName, " ", "_")
		} else {
			file, _ = r.Ffmpeg.DownscaleFile(context, sourceMedia.File)
			fileName = strings.Split(sourceMedia.Filename, ".mp4")[0]
		}

		identifier = fileName + randomString
		fileName = identifier + ".mp4"

		r.Storage.Upload(fileName, file)

		project, _ = r.DB.UpdateProjectSourceMedia(context, database.UpdateProjectSourceMediaParams{
			ID:          project.ID,
			SourceMedia: fileName,
		})

		r.Dubbing.CreateTransformation(context, dubbing.CreateTransformationParams{
			ProjectID: project.ID,
			FileName:  fileName,
			IsSource:  true,
		})

		if initialTargetLanguage != nil {
			r.CreateTranslation(context, project.ID, *initialTargetLanguage, initialLipSync)
		}

	}(newCtx)

	return project, nil
}

// DeleteProject is the resolver for the deleteProject field.
func (r *mutationResolver) DeleteProject(ctx context.Context, projectID int64) (database.Project, error) {
	transformations, _ := r.DB.GetTransformationsByProjectId(ctx, projectID)
	project, _ := r.DB.DeleteProjectById(ctx, projectID)

	newCtx := context.Background()
	go func(ctx context.Context) {
		for _, tfn := range transformations {
			r.Storage.DeleteFile(tfn.TargetMedia)
		}
	}(newCtx)

	return project, nil
}

// CreateTranslation is the resolver for the createTranslation field.
func (r *mutationResolver) CreateTranslation(ctx context.Context, projectID int64, targetLanguage model.SupportedLanguage, lipSync bool) (database.Transformation, error) {
	// fetch source transcript for the project
	sourceTransformation, err := r.DB.GetSourceTransformationByProjectId(ctx, projectID)
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Project Not Processed!")
	}

	// if target transformation already exists, return that
	existingTransformation, err := r.DB.GetTransformationByProjectIdTargetLanguage(ctx, database.GetTransformationByProjectIdTargetLanguageParams{
		ProjectID:      projectID,
		TargetLanguage: targetLanguage.String(),
	})
	if err == nil {
		return existingTransformation, nil
	}

	identifier := fmt.Sprintf("%d-%s-%s", sourceTransformation.ProjectID, utils.GetCurrentDateTimeString(), targetLanguage)
	newFileName := identifier + "_dubbed.mp4"

	// create empty transformation in target language, if target transformation already exists, return that
	newTransformation, _ := r.DB.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      projectID,
		TargetLanguage: targetLanguage.String(),
		TargetMedia:    newFileName,
		Transcript:     pqtype.NullRawMessage{Valid: false, RawMessage: nil},
		IsSource:       false,
		Status:         "starting",
		Progress:       0,
	})

	user := auth.FromContext(ctx)
	newCtx := context.Background()
	newCtx = auth.AttachContext(newCtx, user)

	go func(context context.Context) {
		_, err := r.Dubbing.CreateTranslation(
			newCtx,
			dubbing.CreateTranslationProps{
				SourceTransformation: sourceTransformation,
				TargetTransformation: newTransformation,
				Identifier:           identifier,
				LipSync:              lipSync,
			},
		)

		if err != nil {
			r.Logger.Error("Failed to process transformation", zap.Error(err), zap.Int("project_id", int(projectID)), zap.Int("transformation_id", int(newTransformation.ID)), zap.String("target_language", string(targetLanguage)))
			r.DB.UpdateTransformationStatusById(newCtx, database.UpdateTransformationStatusByIdParams{
				ID:     newTransformation.ID,
				Status: "error",
			})
		}

	}(newCtx)

	return newTransformation, nil
}

// DeleteTransformation is the resolver for the deleteTransformation field.
func (r *mutationResolver) DeleteTransformation(ctx context.Context, transformationID int64) (database.Transformation, error) {
	transformation, _ := r.DB.DeleteTransformationById(ctx, transformationID)

	newCtx := context.Background()
	go func(ctx context.Context) {
		r.Storage.DeleteFile(transformation.TargetMedia)
	}(newCtx)

	return transformation, nil
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
			if t.TargetMedia != "" {
				t.TargetMedia = r.Storage.GetFileLink(t.TargetMedia)
			}
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
		if p.SourceMedia != "" {
			p.SourceMedia = r.Storage.GetFileLink(p.SourceMedia)
		}
		filteredProject = append(filteredProject, p)
	}

	return filteredProject, nil
}

// TargetLanguage is the resolver for the targetLanguage field.
func (r *transformationResolver) TargetLanguage(ctx context.Context, obj *database.Transformation) (model.SupportedLanguage, error) {
	return model.SupportedLanguage(obj.TargetLanguage), nil
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
