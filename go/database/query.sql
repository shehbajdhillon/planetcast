-- name: AddUser :one
INSERT INTO userinfo (email, full_name, created) VALUES ($1, $2, clock_timestamp()) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM userinfo WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM userinfo WHERE id = $1 LIMIT 1;


-- name: CreateTeam :one
INSERT INTO team (slug, name, team_type, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING *;

-- name: GetTeamById :one
SELECT * FROM team WHERE id = $1 LIMIT 1;

-- name: GetTeamBySlug :one
SELECT * FROM team WHERE slug = $1 LIMIT 1;

-- name: AddTeamMembership :one
INSERT INTO team_membership (team_id, user_id, membership_type, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING *;

-- name: GetTeamMemebershipsByUserId :many
SELECT * FROM team_membership WHERE user_id = $1 ORDER BY team_id;

-- name: GetTeamMembershipByTeamIdUserId :one
SELECT * FROM team_membership WHERE team_id = $1 AND user_id = $2 LIMIT 1;


-- name: CreateProject :one
INSERT INTO project (team_id, title, source_language, target_language, source_media, target_media) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetProjectByProjectIdTeamId :one
SELECT * FROM project WHERE id = $1 AND team_id = $2 LIMIT 1;

-- name: GetProjectsByTeamId :many
SELECT * FROM project WHERE team_id = $1;


