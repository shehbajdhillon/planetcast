package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/dubbing"
	"planetcastdev/graph/model"
	"planetcastdev/utils"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	stripe "github.com/stripe/stripe-go/v76"
	portalsession "github.com/stripe/stripe-go/v76/billingportal/session"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/tabbed/pqtype"
	"go.uber.org/zap"
)

// CreateTeam is the resolver for the createTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, teamType database.TeamType, addTrial bool) (database.Team, error) {
	email, _ := auth.EmailFromContext(ctx)
	user, _ := r.DB.GetUserByEmail(ctx, email)

	firstName := strings.ToLower(strings.Split(user.FullName, " ")[0])

	var teamName string
	if teamType == database.TeamTypePERSONAL {
		teamName = fmt.Sprintf("%s's Personal Workspace", strings.Title(firstName))
	} else {
		teamName = fmt.Sprintf("%s's Team", firstName)
	}

	shortUuid := uuid.NewString()[:8]
	teamSlug := fmt.Sprintf("%s-%s", firstName, shortUuid)
	team, err := r.DB.CreateTeam(ctx, database.CreateTeamParams{
		Slug:     teamSlug,
		Name:     teamName,
		TeamType: teamType,
	})

	// Error will likely probably happen if teamSlug is not unique.
	// Although this will obviously not be true all the time
	if err != nil {
		shortUuid = uuid.NewString()[:8]
		teamSlug = fmt.Sprintf("%s-%s", firstName, shortUuid)
		team, err = r.DB.CreateTeam(ctx, database.CreateTeamParams{
			Slug:     teamSlug,
			Name:     teamName,
			TeamType: teamType,
		})
	}

	r.DB.AddTeamMembership(ctx, database.AddTeamMembershipParams{
		TeamID:         team.ID,
		UserID:         user.ID,
		MembershipType: database.MembershipTypeOWNER,
	})

	TRIAL_MINUTES := 0
	if addTrial == true {
		TRIAL_MINUTES = 10
	}

	_, err = r.DB.CreateSubscription(ctx, database.CreateSubscriptionParams{
		TeamID:               team.ID,
		StripeSubscriptionID: sql.NullString{Valid: false, String: ""},
		RemainingCredits:     int64(TRIAL_MINUTES),
	})

	if err != nil {
		r.Logger.Error(
			"Could not add subscription plan for team",
			zap.Error(err),
			zap.Int64("team_id", team.ID),
			zap.String("team_name", team.Name),
		)
	}

	return team, nil
}

// CreateProject is the resolver for the createProject field.
func (r *mutationResolver) CreateProject(ctx context.Context, teamSlug string, title string, sourceMedia *graphql.Upload, youtubeLink *string, uploadOption model.UploadOption, gender string, initialTargetLanguage *string, initialLipSync bool) (database.Project, error) {
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
			r.CreateTranslation(context, project.ID, *initialTargetLanguage, initialLipSync, gender)
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
			if tfn.IsSource == true {
				r.Storage.DeleteFile(fmt.Sprintf("%s-demucs.mp3", tfn.TargetMedia))
			}
		}
	}(newCtx)

	return project, nil
}

// CreateTranslation is the resolver for the createTranslation field.
func (r *mutationResolver) CreateTranslation(ctx context.Context, projectID int64, targetLanguage string, lipSync bool, gender string) (database.Transformation, error) {
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

	var whisperOutput dubbing.WhisperOutput
	json.Unmarshal(sourceTransformation.Transcript.RawMessage, &whisperOutput)
	requiredCredits := r.Dubbing.GetTranscriptLength(&whisperOutput)

	//Check for existing credits, if not enough return
	project, _ := r.DB.GetProjectById(ctx, projectID)
	subs, _ := r.DB.GetSubscriptionsByTeamId(ctx, project.TeamID)
	subPlan := subs[0]

	currentCredits := subPlan.RemainingCredits
	if int64(requiredCredits) > currentCredits {
		return database.Transformation{}, fmt.Errorf("No sufficient credits available to process dubbing. Remaining: %d. Required: %d.", currentCredits, requiredCredits)
	}

	//if enough credits, deduct
	remainingCredits := currentCredits - int64(requiredCredits)
	subPlan, _ = r.DB.SetRemainingCreditsById(ctx, database.SetRemainingCreditsByIdParams{
		RemainingCredits: remainingCredits,
		ID:               subPlan.ID,
	})

	identifier := fmt.Sprintf("%d-%s-%s", sourceTransformation.ProjectID, utils.GetCurrentDateTimeString(), targetLanguage)
	newFileName := identifier + "_dubbed.mp4"

	// create empty transformation in target language, if target transformation already exists, return that
	newTransformation, _ := r.DB.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      projectID,
		TargetLanguage: targetLanguage,
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
				Gender:               gender,
			},
		)

		if err != nil {
			r.Logger.Error("Failed to process transformation", zap.Error(err), zap.Int("project_id", int(projectID)), zap.Int("transformation_id", int(newTransformation.ID)), zap.String("target_language", string(targetLanguage)))
			r.DB.UpdateTransformationStatusById(newCtx, database.UpdateTransformationStatusByIdParams{
				ID:     newTransformation.ID,
				Status: "error",
			})
			//Reimburse credit
			r.DB.SetRemainingCreditsById(newCtx, database.SetRemainingCreditsByIdParams{
				ID:               subPlan.ID,
				RemainingCredits: currentCredits,
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

// CreateCheckoutSession is the resolver for the createCheckoutSession field.
func (r *mutationResolver) CreateCheckoutSession(ctx context.Context, teamSlug string, lookUpKey string) (model.CheckoutSessionResponse, error) {
	production := os.Getenv("PRODUCTION") != ""

	baseUrl := "https://www.planetcast.ai"
	if production == false {
		baseUrl = "http://localhost:3000"
	}

	SuccessURL := fmt.Sprintf("%s/dashboard/%s/settings/subscription?action=success&session_id={CHECKOUT_SESSION_ID}", baseUrl, teamSlug)
	CancelURL := fmt.Sprintf("%s/dashboard/%s/settings/subscription?action=cancel", baseUrl, teamSlug)

	priceLookUpParams := &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{lookUpKey}),
	}

	itemIterable := price.List(priceLookUpParams)

	var price *stripe.Price
	for itemIterable.Next() {
		p := itemIterable.Price()
		price = p
	}

	if price == nil {
		return model.CheckoutSessionResponse{}, fmt.Errorf("No item found")
	}

	var stripeLineItems []*stripe.CheckoutSessionLineItemParams
	stripeLineItems = append(stripeLineItems, &stripe.CheckoutSessionLineItemParams{
		Price:    stripe.String(price.ID),
		Quantity: stripe.Int64(1),
	})

	currentCustomer, err := r.Payments.GetCustomerByTeamSlug(ctx, teamSlug)

	params := &stripe.CheckoutSessionParams{
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          stripeLineItems,
		SuccessURL:         stripe.String(SuccessURL),
		CancelURL:          stripe.String(CancelURL),
		Customer:           stripe.String(currentCustomer.ID),
	}

	session, err := session.New(params)

	if err != nil {
		return model.CheckoutSessionResponse{}, err
	}

	return model.CheckoutSessionResponse{SessionID: session.ID}, nil
}

// CreatePortalSession is the resolver for the createPortalSession field.
func (r *mutationResolver) CreatePortalSession(ctx context.Context, teamSlug string) (model.PortalSessionResponse, error) {
	production := os.Getenv("PRODUCTION") != ""

	baseUrl := "https://www.planetcast.ai"
	if production == false {
		baseUrl = "http://localhost:3000"
	}

	currentTeam, err := r.DB.GetTeamBySlug(ctx, teamSlug)
	if err != nil || currentTeam.StripeCustomerID.Valid == false {
		return model.PortalSessionResponse{}, fmt.Errorf("No customer records found")
	}

	ReturnURL := fmt.Sprintf("%s/dashboard/%s/settings/subscription", baseUrl, teamSlug)
	customerId := currentTeam.StripeCustomerID.String
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerId),
		ReturnURL: stripe.String(ReturnURL),
	}

	ps, err := portalsession.New(params)
	if err != nil {
		return model.PortalSessionResponse{}, err
	}

	return model.PortalSessionResponse{SessionURL: ps.URL}, nil
}

// SendTeamInvite is the resolver for the sendTeamInvite field.
func (r *mutationResolver) SendTeamInvite(ctx context.Context, teamSlug string, inviteeEmail string) (bool, error) {
	team, _ := r.DB.GetTeamBySlug(ctx, teamSlug)
	user, err := r.DB.GetUserByEmail(ctx, inviteeEmail)
	if err == nil {
		_, err = r.DB.GetTeamMembershipByTeamIdUserId(ctx, database.GetTeamMembershipByTeamIdUserIdParams{
			TeamID: team.ID,
			UserID: user.ID,
		})

		if err == nil {
			return false, fmt.Errorf("User already a member of the team")
		}
	}

	_, err = r.DB.AddTeamInvite(ctx, database.AddTeamInviteParams{
		Slug:         uuid.NewString(),
		InviteeEmail: inviteeEmail,
		TeamID:       team.ID,
	})

	if err != nil {
		return false, fmt.Errorf("User already invited to the team")
	}

	return true, nil
}

// DeleteTeamInvite is the resolver for the deleteTeamInvite field.
func (r *mutationResolver) DeleteTeamInvite(ctx context.Context, inviteSlug string) (bool, error) {
	_, err := r.DB.DeleteTeamInviteBySlug(ctx, inviteSlug)
	if err != nil {
		return false, fmt.Errorf("Could not remove invite")
	}
	return true, nil
}

// AcceptTeamInvite is the resolver for the acceptTeamInvite field.
func (r *mutationResolver) AcceptTeamInvite(ctx context.Context, inviteSlug string) (bool, error) {
	invite, _ := r.DB.GetTeamInviteBySlug(ctx, inviteSlug)
	user, _ := r.DB.GetUserByEmail(ctx, invite.InviteeEmail)
	_, err := r.DB.AddTeamMembership(ctx, database.AddTeamMembershipParams{
		TeamID:         invite.TeamID,
		UserID:         user.ID,
		MembershipType: database.MembershipTypeMEMBER,
	})
	if err != nil {
		return false, fmt.Errorf("Could not process invite")
	}
	r.DB.DeleteTeamInviteBySlug(ctx, inviteSlug)
	return true, nil
}

// DubbingCreditsRequired is the resolver for the dubbingCreditsRequired field.
func (r *projectResolver) DubbingCreditsRequired(ctx context.Context, obj *database.Project) (*int64, error) {
	sourceTransformation, err := r.DB.GetSourceTransformationByProjectId(ctx, obj.ID)
	if err != nil {
		return nil, nil
	}
	var whisperOutput dubbing.WhisperOutput
	json.Unmarshal(sourceTransformation.Transcript.RawMessage, &whisperOutput)
	requiredCredits := int64(r.Dubbing.GetTranscriptLength(&whisperOutput))

	return &requiredCredits, nil
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
	memberships, _ := r.DB.GetTeamMembershipsByUserId(ctx, user.ID)
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

// GetUserInfo is the resolver for the getUserInfo field.
func (r *queryResolver) GetUserInfo(ctx context.Context) (model.AccountInfo, error) {
	email, _ := auth.EmailFromContext(ctx)
	user, _ := r.DB.GetUserByEmail(ctx, email)

	memberships, _ := r.DB.GetTeamMembershipsByUserId(ctx, user.ID)
	invites, _ := r.DB.GetTeamInvitesByInviteeEmail(ctx, user.Email)

	return model.AccountInfo{User: user, Invites: invites, Teams: memberships}, nil
}

// StripeSubscriptionID is the resolver for the stripeSubscriptionId field.
func (r *subscriptionPlanResolver) StripeSubscriptionID(ctx context.Context, obj *database.SubscriptionPlan) (*string, error) {
	if obj.StripeSubscriptionID.Valid == false {
		return nil, nil
	}
	subscriptionId := obj.StripeSubscriptionID.String
	return &subscriptionId, nil
}

// SubscriptionData is the resolver for the subscriptionData field.
func (r *subscriptionPlanResolver) SubscriptionData(ctx context.Context, obj *database.SubscriptionPlan) (*model.SubscriptionData, error) {
	if obj.StripeSubscriptionID.Valid == false {
		return nil, nil
	}
	subscriptionData, err := r.Payments.GetSubscriptionPlanData(obj.StripeSubscriptionID.String)
	if err != nil {
		r.Logger.Error("Could not fetch subscription plan data %s", zap.Error(err))
		return nil, fmt.Errorf("Could not fetch subscription plan data %s", err.Error())
	}
	return subscriptionData, nil
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

// SubscriptionPlans is the resolver for the subscriptionPlans field.
func (r *teamResolver) SubscriptionPlans(ctx context.Context, obj *database.Team, subscriptionID *int64) ([]database.SubscriptionPlan, error) {
	if subscriptionID == nil {
		return r.DB.GetSubscriptionsByTeamId(ctx, obj.ID)
	}

	subscription, err := r.DB.GetSubscriptionByTeamIdSubscriptionId(ctx, database.GetSubscriptionByTeamIdSubscriptionIdParams{TeamID: obj.ID, ID: *subscriptionID})
	return []database.SubscriptionPlan{subscription}, err
}

// Members is the resolver for the members field.
func (r *teamResolver) Members(ctx context.Context, obj *database.Team) ([]database.TeamMembership, error) {
	memberships, _ := r.DB.GetTeamMembershipsByTeamId(ctx, obj.ID)
	return memberships, nil
}

// Invitees is the resolver for the invitees field.
func (r *teamResolver) Invitees(ctx context.Context, obj *database.Team) ([]database.TeamInvite, error) {
	invites, err := r.DB.GetTeamInvitesByTeamId(ctx, obj.ID)
	if err != nil {
		return []database.TeamInvite{}, nil
	}
	inviteeEmails := []database.TeamInvite{}
	for _, invite := range invites {
		inviteeEmails = append(inviteeEmails, database.TeamInvite{
			Slug:         invite.Slug,
			InviteeEmail: invite.InviteeEmail,
		})
	}
	return inviteeEmails, nil
}

// InviteSlug is the resolver for the inviteSlug field.
func (r *teamInviteResolver) InviteSlug(ctx context.Context, obj *database.TeamInvite) (string, error) {
	return obj.Slug, nil
}

// TeamName is the resolver for the teamName field.
func (r *teamInviteResolver) TeamName(ctx context.Context, obj *database.TeamInvite) (string, error) {
	team, _ := r.DB.GetTeamById(ctx, obj.TeamID)
	return team.Name, nil
}

// MembershipType is the resolver for the membershipType field.
func (r *teamMembershipResolver) MembershipType(ctx context.Context, obj *database.TeamMembership) (string, error) {
	return string(obj.MembershipType), nil
}

// User is the resolver for the user field.
func (r *teamMembershipResolver) User(ctx context.Context, obj *database.TeamMembership) (database.Userinfo, error) {
	return r.DB.GetUserById(ctx, obj.UserID)
}

// TeamSlug is the resolver for the teamSlug field.
func (r *teamMembershipResolver) TeamSlug(ctx context.Context, obj *database.TeamMembership) (string, error) {
	team, _ := r.DB.GetTeamById(ctx, obj.TeamID)
	return team.Slug, nil
}

// TeamName is the resolver for the teamName field.
func (r *teamMembershipResolver) TeamName(ctx context.Context, obj *database.TeamMembership) (string, error) {
	team, _ := r.DB.GetTeamById(ctx, obj.TeamID)
	return team.Name, nil
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

// SubscriptionPlan returns SubscriptionPlanResolver implementation.
func (r *Resolver) SubscriptionPlan() SubscriptionPlanResolver { return &subscriptionPlanResolver{r} }

// Team returns TeamResolver implementation.
func (r *Resolver) Team() TeamResolver { return &teamResolver{r} }

// TeamInvite returns TeamInviteResolver implementation.
func (r *Resolver) TeamInvite() TeamInviteResolver { return &teamInviteResolver{r} }

// TeamMembership returns TeamMembershipResolver implementation.
func (r *Resolver) TeamMembership() TeamMembershipResolver { return &teamMembershipResolver{r} }

// Transformation returns TransformationResolver implementation.
func (r *Resolver) Transformation() TransformationResolver { return &transformationResolver{r} }

type mutationResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionPlanResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
type teamInviteResolver struct{ *Resolver }
type teamMembershipResolver struct{ *Resolver }
type transformationResolver struct{ *Resolver }
