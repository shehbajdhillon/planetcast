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

-- name: UpdateTeamStripeCustomerIdByTeamId :one
UPDATE team SET stripe_customer_id = $2 WHERE id = $1 RETURNING *;

-- name: GetTeamByStripeCustomerId :one
SELECT * FROM team WHERE stripe_customer_id = $1;


-- name: AddTeamMembership :one
INSERT INTO team_membership (team_id, user_id, membership_type, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING *;

-- name: GetTeamMemebershipsByUserId :many
SELECT * FROM team_membership WHERE user_id = $1 ORDER BY team_id;

-- name: GetTeamMembershipByTeamIdUserId :one
SELECT * FROM team_membership WHERE team_id = $1 AND user_id = $2 LIMIT 1;


-- name: CreateSubscription :one
INSERT INTO subscription_plan
(team_id, stripe_subscription_id, remaining_credits, created)
VALUES ($1, $2, $3, clock_timestamp()) RETURNING *;


-- name: GetSubscriptionsByTeamId :many
SELECT * FROM subscription_plan WHERE team_id = $1 ORDER BY created;

-- name: GetSubscriptionByTeamIdSubscriptionId :one
SELECT * FROM subscription_plan WHERE team_id = $1 AND id = $2 LIMIT 1;

-- name: GetSubscriptionById :one
SELECT * FROM subscription_plan WHERE id = $1 LIMIT 1;

-- name: GetSubscriptionByStripeSubscriptionId :one
SELECT * FROM subscription_plan WHERE stripe_subscription_id = $1 LIMIT 1;

-- name: AddSubscriptionCreditsByTeamId :one
UPDATE subscription_plan SET remaining_credits = remaining_credits + $2 WHERE team_id = $1 RETURNING *;

-- name: SetSubscriptionStripeIdByTeamId :one
UPDATE subscription_plan SET stripe_subscription_id = $2 WHERE team_id = $1 RETURNING *;

-- name: SetRemainingCreditsById :one
UPDATE subscription_plan SET remaining_credits = $2 WHERE id = $1 RETURNING *;


-- name: CreateProject :one
INSERT INTO project (team_id, title, source_media, created) VALUES ($1, $2, $3, clock_timestamp()) RETURNING *;

-- name: GetProjectById :one
SELECT * FROM project WHERE id = $1 LIMIT 1;

-- name: GetProjectByProjectIdTeamId :one
SELECT * FROM project WHERE id = $1 AND team_id = $2 LIMIT 1;

-- name: GetProjectsByTeamId :many
SELECT * FROM project WHERE team_id = $1 ORDER BY created;

-- name: UpdateProjectSourceMedia :one
UPDATE project SET source_media = $2 WHERE id = $1 RETURNING *;

-- name: DeleteProjectById :one
DELETE FROM project WHERE id = $1 RETURNING *;


-- name: CreateTransformation :one
INSERT INTO transformation
(project_id, target_language, target_media, transcript, is_source, status, progress, created)
VALUES ($1, $2, $3, $4, $5, $6, $7, clock_timestamp()) RETURNING *;

-- name: UpdateTranscriptById :one
UPDATE transformation SET transcript = $2 WHERE id = $1 RETURNING *;

-- name: UpdateTargetMediaById :one
UPDATE transformation SET target_media = $2 WHERE id = $1 RETURNING *;

-- name: UpdateTransformationStatusById :one
UPDATE transformation SET status = $2 WHERE id = $1 RETURNING *;

-- name: UpdateTransformationProgressById :one
UPDATE transformation SET progress = $2 WHERE id = $1 RETURNING *;

-- name: GetTransformationById :one
SELECT * FROM transformation WHERE id = $1 LIMIT 1;

-- name: GetTransformationsByProjectId :many
SELECT * FROM transformation WHERE project_id = $1 ORDER BY created;

-- name: GetTransformationByTransformationIdProjectId :one
SELECT * FROM transformation WHERE id = $1 AND project_id = $2 LIMIT 1;

-- name: GetSourceTransformationByProjectId :one
SELECT * FROM transformation WHERE project_id = $1 AND is_source = true LIMIT 1;

-- name: GetTransformationByProjectIdTargetLanguage :one
SELECT * FROM transformation WHERE project_id = $1 AND target_language = $2 LIMIT 1;

-- name: DeleteTransformationById :one
DELETE FROM transformation WHERE id = $1 RETURNING *;

