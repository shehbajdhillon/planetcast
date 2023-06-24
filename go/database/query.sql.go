// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: query.sql

package database

import (
	"context"
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
INSERT INTO project (team_id, title, source_language, target_language, source_media, target_media) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, team_id, title, source_language, target_language, source_media, target_media
`

type CreateProjectParams struct {
	TeamID         int64
	Title          string
	SourceLanguage SupportedLanguage
	TargetLanguage SupportedLanguage
	SourceMedia    string
	TargetMedia    string
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.TeamID,
		arg.Title,
		arg.SourceLanguage,
		arg.TargetLanguage,
		arg.SourceMedia,
		arg.TargetMedia,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.TargetLanguage,
		&i.SourceMedia,
		&i.TargetMedia,
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

const deleteProjectById = `-- name: DeleteProjectById :one
DELETE FROM project WHERE id = $1 RETURNING id, team_id, title, source_language, target_language, source_media, target_media
`

func (q *Queries) DeleteProjectById(ctx context.Context, id int64) (Project, error) {
	row := q.db.QueryRowContext(ctx, deleteProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.TargetLanguage,
		&i.SourceMedia,
		&i.TargetMedia,
	)
	return i, err
}

const getProjectById = `-- name: GetProjectById :one
SELECT id, team_id, title, source_language, target_language, source_media, target_media FROM project WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProjectById(ctx context.Context, id int64) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.Title,
		&i.SourceLanguage,
		&i.TargetLanguage,
		&i.SourceMedia,
		&i.TargetMedia,
	)
	return i, err
}

const getProjectByProjectIdTeamId = `-- name: GetProjectByProjectIdTeamId :one
SELECT id, team_id, title, source_language, target_language, source_media, target_media FROM project WHERE id = $1 AND team_id = $2 LIMIT 1
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
		&i.TargetLanguage,
		&i.SourceMedia,
		&i.TargetMedia,
	)
	return i, err
}

const getProjectsByTeamId = `-- name: GetProjectsByTeamId :many
SELECT id, team_id, title, source_language, target_language, source_media, target_media FROM project WHERE team_id = $1
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
			&i.TargetLanguage,
			&i.SourceMedia,
			&i.TargetMedia,
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
