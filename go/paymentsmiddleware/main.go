package paymentsmiddleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.uber.org/zap"
)

type Payments struct {
	secretKey string
	database  *database.Queries
	logger    *zap.Logger
}

type PaymentsConnectProps struct {
	Logger   *zap.Logger
	Database *database.Queries
}

func Connect(args PaymentsConnectProps) *Payments {
	STRIPE_KEY := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = STRIPE_KEY
	return &Payments{secretKey: STRIPE_KEY, database: args.Database, logger: args.Logger}
}

// Customer Management
func (p *Payments) createCustomer(email string, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	return customer.New(params)
}

func (p *Payments) getCustomer(customerId string) (*stripe.Customer, error) {
	return customer.Get(customerId, nil)
}

func (p *Payments) UpdateCustomer(customerId string, newEmail string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(newEmail),
	}
	return customer.Update(customerId, params)
}

func (p *Payments) GetCustomerByTeamSlug(ctx context.Context, teamSlug string) (*stripe.Customer, error) {
	team, err := p.database.GetTeamBySlug(ctx, teamSlug)
	if err != nil {
		p.logger.Error("Unable to fetch team from DB", zap.Error(err), zap.Int64("team_id", team.ID), zap.String("team_name", team.Name))
		return nil, err
	}

	if team.StripeCustomerID.Valid == true {
		return p.getCustomer(team.StripeCustomerID.String)
	}

	userEmail, err := auth.EmailFromContext(ctx)
	if err != nil {
		p.logger.Error("Unable to fetch user email to assign to stripe customer object")
		return nil, err
	}

	newCustomer, err := p.createCustomer(userEmail, team.Name)
	if err != nil {
		p.logger.Error("Unable to create stripe customer for team", zap.Error(err), zap.Int64("team_id", team.ID), zap.String("team_name", team.Name))
		return nil, err
	}

	_, err = p.database.UpdateTeamStripeCustomerIdByTeamId(ctx, database.UpdateTeamStripeCustomerIdByTeamIdParams{
		ID:               team.ID,
		StripeCustomerID: sql.NullString{Valid: true, String: newCustomer.ID},
	})
	if err != nil {
		p.logger.Error("Unable to update stripe customer for team in DB", zap.Error(err), zap.Int64("team_id", team.ID), zap.String("team_name", team.Name))
		return nil, err
	}

	return newCustomer, nil
}

func (p *Payments) DeleteCustomer(customerId string) (*stripe.Customer, error) {
	return customer.Del(customerId, nil)
}

// Subscription Management
func (p *Payments) CreateSubcription(customerId string, priceId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerId),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(priceId)},
		},
	}
	return subscription.New(params)
}

func (p *Payments) GetSubcription(subscriptionId string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionId, nil)
}

func (p *Payments) CancelSubcription(subscriptionId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{
		InvoiceNow: stripe.Bool(true),
		Prorate:    stripe.Bool(true),
	}
	return subscription.Cancel(subscriptionId, params)
}

func (p *Payments) ListSubscriptions(customerId string) ([]*stripe.Subscription, error) {
	params := &stripe.SubscriptionListParams{}
	params.Filters.AddFilter("customer", "", customerId)
	i := subscription.List(params)
	var subscriptions []*stripe.Subscription
	for i.Next() {
		subscriptions = append(subscriptions, i.Subscription())
	}

	return subscriptions, i.Err()
}

// Payment Management
func (p *Payments) AttachPaymentMethod(customerId string, paymentMethodId string) (*stripe.PaymentMethod, error) {
	return paymentmethod.Detach(paymentMethodId, nil)
}

func (p *Payments) DetachPaymentMethod(paymentMethodId string) (*stripe.PaymentMethod, error) {
	return paymentmethod.Detach(paymentMethodId, nil)
}

func ListPaymentMethods(customerId string, paymentMethodType string) ([]*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerId),
		Type:     stripe.String(paymentMethodType),
	}

	i := paymentmethod.List(params)

	var paymentMethods []*stripe.PaymentMethod
	for i.Next() {
		paymentMethods = append(paymentMethods, i.PaymentMethod())
	}

	return paymentMethods, i.Err()
}

// Webhook Handler
func HandleStripeWebhook(w http.ResponseWriter, r *http.Request, endpointSecret string) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error parsing webhook JSON")
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	_, err = webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error verifying webhook signature")
		return
	}

	switch event.Type {

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error parsing payment intent")
			return
		}

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error parsing payment intent")
			return
		}

	case "customer.subcription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			RespondWithJSON(w, http.StatusInternalServerError, "Error parsing subscription")
			return
		}

	default:
		RespondWithError(w, http.StatusNotFound, "Event type not handled")
		return
	}

	RespondWithError(w, http.StatusNotFound, "Event type not handled")
	return

}

// RespondWithError sends an error response with a given HTTP status code and error message.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// FormatStripeAmount formats amount in cents to a float value.
func FormatStripeAmount(amount int64) float64 {
	return float64(amount) / 100
}
