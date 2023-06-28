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
INSERT INTO project (team_id, title, source_language, source_media) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetProjectById :one
SELECT * FROM project WHERE id = $1 LIMIT 1;

-- name: GetProjectByProjectIdTeamId :one
SELECT * FROM project WHERE id = $1 AND team_id = $2 LIMIT 1;

-- name: GetProjectsByTeamId :many
SELECT * FROM project WHERE team_id = $1;

-- name: DeleteProjectById :one
DELETE FROM project WHERE id = $1 RETURNING *;


-- name: CreateTransformation :one
INSERT INTO transformation (project_id, target_language, target_media, transcript) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateTranscriptById :one
UPDATE transformation SET transcript = $2 WHERE project_id = $1 RETURNING *;

-- name: UpdateTargetMediaById :one
UPDATE transformation SET target_media = $2 WHERE project_id = $1 RETURNING *;

-- name: GetTransformationById :one
SELECT * FROM transformation WHERE id = $1 LIMIT 1;

-- name: GetTransformationsByProjectId :many
SELECT * FROM transformation WHERE project_id = $1;

-- name: GetTransformationByTransformationIdProjectId :one
SELECT * FROM transformation WHERE id = $1 AND project_id = $2 LIMIT 1;
