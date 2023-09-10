// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: query.sql

package database

import (
	"context"

	"github.com/tabbed/pqtype"
)

const addTeamMembership = `-- name: AddTeamMembership :one
INSERT INTO team_membership (team_id, user_id, membership_type, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING id, team_id, user_id, membership_type, created
`

type AddTeamMembershipParams struct {
	TeamID         int64
	UserID         int64
	MembershipType MembershipType
}

func (q *Queries) AddTeamMembership(ctx context.Context, arg AddTeamMembershipParams) (TeamMembership, error) {
	row := q.db.QueryRowContext(ctx, addTeamMembership, arg.TeamID, arg.UserID, arg.MembershipType)
	var i TeamMembership
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.UserID,
		&i.MembershipType,
		&i.Created,
	)
	return i, err
}

const addUser = `-- name: AddUser :one
INSERT INTO userinfo (email, full_name, created) VALUES ($1, $2, clock_timestamp()) RETURNING id, email, full_name, created
`

type AddUserParams struct {
	Email    string
	FullName string
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (Userinfo, error) {
	row := q.db.QueryRowContext(ctx, addUser, arg.Email, arg.FullName)
	var i Userinfo
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Created,
	)
	return i, err
}

const createProject = `-- name: CreateProject :one
INSERT INTO project (team_id, title, source_language, source_media) VALUES ($1, $2, $3, $4) RETURNING id, team_id, title, source_language, source_media
`

type CreateProjectParams struct {
	TeamID         int64
	Title          string
	SourceLanguage string
	SourceMedia    string
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.TeamID,
		arg.Title,
		arg.SourceLanguage,
		arg.SourceMedia,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.SourceMedia,
	)
	return i, err
}

const createTeam = `-- name: CreateTeam :one
INSERT INTO team (slug, name, team_type, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING id, slug, name, team_type, created
`

type CreateTeamParams struct {
	Slug     string
	Name     string
	TeamType TeamType
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, createTeam, arg.Slug, arg.Name, arg.TeamType)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Name,
		&i.TeamType,
		&i.Created,
	)
	return i, err
}

const createTransformation = `-- name: CreateTransformation :one
INSERT INTO transformation (project_id, target_language, target_media, transcript, is_source, status, progress) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, project_id, target_language, target_media, transcript, is_source, status, progress
`

type CreateTransformationParams struct {
	ProjectID      int64
	TargetLanguage string
	TargetMedia    string
	Transcript     pqtype.NullRawMessage
	IsSource       bool
	Status         string
	Progress       float64
}

func (q *Queries) CreateTransformation(ctx context.Context, arg CreateTransformationParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, createTransformation,
		arg.ProjectID,
		arg.TargetLanguage,
		arg.TargetMedia,
		arg.Transcript,
		arg.IsSource,
		arg.Status,
		arg.Progress,
	)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const deleteProjectById = `-- name: DeleteProjectById :one
DELETE FROM project WHERE id = $1 RETURNING id, team_id, title, source_language, source_media
`

func (q *Queries) DeleteProjectById(ctx context.Context, id int64) (Project, error) {
	row := q.db.QueryRowContext(ctx, deleteProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.SourceMedia,
	)
	return i, err
}

const getProjectById = `-- name: GetProjectById :one
SELECT id, team_id, title, source_language, source_media FROM project WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProjectById(ctx context.Context, id int64) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.SourceMedia,
	)
	return i, err
}

const getProjectByProjectIdTeamId = `-- name: GetProjectByProjectIdTeamId :one
SELECT id, team_id, title, source_language, source_media FROM project WHERE id = $1 AND team_id = $2 LIMIT 1
`

type GetProjectByProjectIdTeamIdParams struct {
	ID     int64
	TeamID int64
}

func (q *Queries) GetProjectByProjectIdTeamId(ctx context.Context, arg GetProjectByProjectIdTeamIdParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProjectByProjectIdTeamId, arg.ID, arg.TeamID)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.SourceMedia,
	)
	return i, err
}

const getProjectsByTeamId = `-- name: GetProjectsByTeamId :many
SELECT id, team_id, title, source_language, source_media FROM project WHERE team_id = $1
`

func (q *Queries) GetProjectsByTeamId(ctx context.Context, teamID int64) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, getProjectsByTeamId, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.Title,
			&i.SourceLanguage,
			&i.SourceMedia,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSourceTransformationByProjectId = `-- name: GetSourceTransformationByProjectId :one
SELECT id, project_id, target_language, target_media, transcript, is_source, status, progress FROM transformation WHERE project_id = $1 AND is_source = true LIMIT 1
`

func (q *Queries) GetSourceTransformationByProjectId(ctx context.Context, projectID int64) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, getSourceTransformationByProjectId, projectID)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const getTeamById = `-- name: GetTeamById :one
SELECT id, slug, name, team_type, created FROM team WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTeamById(ctx context.Context, id int64) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeamById, id)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Name,
		&i.TeamType,
		&i.Created,
	)
	return i, err
}

const getTeamBySlug = `-- name: GetTeamBySlug :one
SELECT id, slug, name, team_type, created FROM team WHERE slug = $1 LIMIT 1
`

func (q *Queries) GetTeamBySlug(ctx context.Context, slug string) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeamBySlug, slug)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Name,
		&i.TeamType,
		&i.Created,
	)
	return i, err
}

const getTeamMembershipByTeamIdUserId = `-- name: GetTeamMembershipByTeamIdUserId :one
SELECT id, team_id, user_id, membership_type, created FROM team_membership WHERE team_id = $1 AND user_id = $2 LIMIT 1
`

type GetTeamMembershipByTeamIdUserIdParams struct {
	TeamID int64
	UserID int64
}

func (q *Queries) GetTeamMembershipByTeamIdUserId(ctx context.Context, arg GetTeamMembershipByTeamIdUserIdParams) (TeamMembership, error) {
	row := q.db.QueryRowContext(ctx, getTeamMembershipByTeamIdUserId, arg.TeamID, arg.UserID)
	var i TeamMembership
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.UserID,
		&i.MembershipType,
		&i.Created,
	)
	return i, err
}

const getTeamMemebershipsByUserId = `-- name: GetTeamMemebershipsByUserId :many
SELECT id, team_id, user_id, membership_type, created FROM team_membership WHERE user_id = $1 ORDER BY team_id
`

func (q *Queries) GetTeamMemebershipsByUserId(ctx context.Context, userID int64) ([]TeamMembership, error) {
	rows, err := q.db.QueryContext(ctx, getTeamMemebershipsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TeamMembership
	for rows.Next() {
		var i TeamMembership
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.UserID,
			&i.MembershipType,
			&i.Created,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransformationById = `-- name: GetTransformationById :one
SELECT id, project_id, target_language, target_media, transcript, is_source, status, progress FROM transformation WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransformationById(ctx context.Context, id int64) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, getTransformationById, id)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const getTransformationByProjectIdTargetLanguage = `-- name: GetTransformationByProjectIdTargetLanguage :one
SELECT id, project_id, target_language, target_media, transcript, is_source, status, progress FROM transformation WHERE project_id = $1 AND target_language = $2 LIMIT 1
`

type GetTransformationByProjectIdTargetLanguageParams struct {
	ProjectID      int64
	TargetLanguage string
}

func (q *Queries) GetTransformationByProjectIdTargetLanguage(ctx context.Context, arg GetTransformationByProjectIdTargetLanguageParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, getTransformationByProjectIdTargetLanguage, arg.ProjectID, arg.TargetLanguage)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const getTransformationByTransformationIdProjectId = `-- name: GetTransformationByTransformationIdProjectId :one
SELECT id, project_id, target_language, target_media, transcript, is_source, status, progress FROM transformation WHERE id = $1 AND project_id = $2 LIMIT 1
`

type GetTransformationByTransformationIdProjectIdParams struct {
	ID        int64
	ProjectID int64
}

func (q *Queries) GetTransformationByTransformationIdProjectId(ctx context.Context, arg GetTransformationByTransformationIdProjectIdParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, getTransformationByTransformationIdProjectId, arg.ID, arg.ProjectID)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const getTransformationsByProjectId = `-- name: GetTransformationsByProjectId :many
SELECT id, project_id, target_language, target_media, transcript, is_source, status, progress FROM transformation WHERE project_id = $1
`

func (q *Queries) GetTransformationsByProjectId(ctx context.Context, projectID int64) ([]Transformation, error) {
	rows, err := q.db.QueryContext(ctx, getTransformationsByProjectId, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transformation
	for rows.Next() {
		var i Transformation
		if err := rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.TargetLanguage,
			&i.TargetMedia,
			&i.Transcript,
			&i.IsSource,
			&i.Status,
			&i.Progress,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, full_name, created FROM userinfo WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (Userinfo, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i Userinfo
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Created,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, full_name, created FROM userinfo WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (Userinfo, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i Userinfo
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Created,
	)
	return i, err
}

const updateTargetMediaById = `-- name: UpdateTargetMediaById :one
UPDATE transformation SET target_media = $2 WHERE id = $1 RETURNING id, project_id, target_language, target_media, transcript, is_source, status, progress
`

type UpdateTargetMediaByIdParams struct {
	ID          int64
	TargetMedia string
}

func (q *Queries) UpdateTargetMediaById(ctx context.Context, arg UpdateTargetMediaByIdParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, updateTargetMediaById, arg.ID, arg.TargetMedia)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const updateTranscriptById = `-- name: UpdateTranscriptById :one
UPDATE transformation SET transcript = $2 WHERE id = $1 RETURNING id, project_id, target_language, target_media, transcript, is_source, status, progress
`

type UpdateTranscriptByIdParams struct {
	ID         int64
	Transcript pqtype.NullRawMessage
}

func (q *Queries) UpdateTranscriptById(ctx context.Context, arg UpdateTranscriptByIdParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, updateTranscriptById, arg.ID, arg.Transcript)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const updateTransformationProgressById = `-- name: UpdateTransformationProgressById :one
UPDATE transformation SET progress = $2 WHERE id = $1 RETURNING id, project_id, target_language, target_media, transcript, is_source, status, progress
`

type UpdateTransformationProgressByIdParams struct {
	ID       int64
	Progress float64
}

func (q *Queries) UpdateTransformationProgressById(ctx context.Context, arg UpdateTransformationProgressByIdParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, updateTransformationProgressById, arg.ID, arg.Progress)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}

const updateTransformationStatusById = `-- name: UpdateTransformationStatusById :one
UPDATE transformation SET status = $2 WHERE id = $1 RETURNING id, project_id, target_language, target_media, transcript, is_source, status, progress
`

type UpdateTransformationStatusByIdParams struct {
	ID     int64
	Status string
}

func (q *Queries) UpdateTransformationStatusById(ctx context.Context, arg UpdateTransformationStatusByIdParams) (Transformation, error) {
	row := q.db.QueryRowContext(ctx, updateTransformationStatusById, arg.ID, arg.Status)
	var i Transformation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.TargetLanguage,
		&i.TargetMedia,
		&i.Transcript,
		&i.IsSource,
		&i.Status,
		&i.Progress,
	)
	return i, err
}
