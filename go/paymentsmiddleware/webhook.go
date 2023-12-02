package paymentsmiddleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"planetcastdev/database"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.uber.org/zap"
)

func (p *Payments) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	bodyReader := http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(bodyReader)
	if err != nil {
		p.RespondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	ctx := context.Background()

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		p.RespondWithError(w, http.StatusBadRequest, "Error parsing webhook JSON")
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_ENDPOINT_SECRET")
	sigHeader := r.Header.Get("Stripe-Signature")
	_, err = webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		p.RespondWithError(w, http.StatusBadRequest, "Error verifying webhook signature")
		return
	}

	switch event.Type {

	case "customer.deleted":
		customer, err := p.parseCustomerBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing customer data "+err.Error())
			return
		}
		err = p.DeleteCustomerFromDB(ctx, customer.ID)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling customer deletion "+err.Error())
			return
		}

	case "customer.subscription.deleted":
		subscription, err := p.parseSubscriptionBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing subscription data "+err.Error())
			return
		}
		err = p.handleSubscriptionDeleted(ctx, *subscription)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling subscription deletion "+err.Error())
			return
		}

	case "customer.subscription.updated":
		subscription, err := p.parseSubscriptionBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing subscription data "+err.Error())
			return
		}
		err = p.handleSubscriptionUpdated(ctx, *subscription)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling subscription update "+err.Error())
			return
		}

	case "customer.subscription.created":
		subscription, err := p.parseSubscriptionBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing subscription data "+err.Error())
			return
		}
		err = p.handleSubscriptionCreated(ctx, *subscription)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling subscription creation "+err.Error())
			return
		}

	case "invoice.paid":
		invoice, err := p.parseInvoiceBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing invoice data "+err.Error())
			return
		}
		err = p.handleInvoicePaid(ctx, *invoice)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling invoide paid event "+err.Error())
			return
		}

	case "invoice.payment_failed":
		invoice, err := p.parseInvoiceBody(event.Data.Raw)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error parsing invoice data "+err.Error())
			return
		}
		err = p.handleInvoicePaymentFailed(ctx, *invoice)
		if err != nil {
			p.RespondWithError(w, http.StatusInternalServerError, "Error handling invoide payment failed event "+err.Error())
			return
		}

	default:
		p.RespondWithError(w, http.StatusNotFound, "Event type not handled: "+string(event.Type))
		return
	}
}

func (p *Payments) parseCustomerBody(jsonMessage json.RawMessage) (*stripe.Customer, error) {
	var customer stripe.Customer
	err := json.Unmarshal(jsonMessage, &customer)
	if err != nil {
		return nil, fmt.Errorf("Could not parse customer body data: %s", err.Error())
	}
	return &customer, nil
}

func (p *Payments) parseInvoiceBody(jsonMessage json.RawMessage) (*stripe.Invoice, error) {
	var invoice stripe.Invoice
	err := json.Unmarshal(jsonMessage, &invoice)
	if err != nil {
		return nil, fmt.Errorf("Could not parse invoice body data: %s", err.Error())
	}
	return &invoice, nil
}

func (p *Payments) parseSubscriptionBody(jsonMessage json.RawMessage) (*stripe.Subscription, error) {
	var subscription stripe.Subscription
	err := json.Unmarshal(jsonMessage, &subscription)
	if err != nil {
		return nil, fmt.Errorf("Could not parse subscription data: %s", err.Error())
	}
	return &subscription, nil
}

func (p *Payments) handleSubscriptionDeleted(ctx context.Context, subscription stripe.Subscription) error {
	customerId := subscription.Customer.ID
	p.logCustomer(customerId, "Customer Subscription Deleted")

	sub_plan, err := p.database.GetSubscriptionByStripeSubscriptionId(ctx, sql.NullString{Valid: true, String: subscription.ID})
	if err != nil {
		return fmt.Errorf("Could not find active subscription with id %s: %s", subscription.ID, err.Error())
	}

	teamId := sub_plan.TeamID

	team, err := p.database.GetTeamById(ctx, teamId)
	if err != nil {
		return fmt.Errorf("Could not fetch team from DB using stripe customer ID %s: %s", customerId, err.Error())
	}

	_, err = p.database.SetSubscriptionStripeIdByTeamId(ctx, database.SetSubscriptionStripeIdByTeamIdParams{
		TeamID:               teamId,
		StripeSubscriptionID: sql.NullString{Valid: false, String: ""},
	})

	if err != nil {
		return fmt.Errorf("Unable to update stripe subscription id for team: %d: %s", teamId, err.Error())
	}

	p.logger.Info(
		"Set subscription to unactive due to subscription deletion by team",
		zap.String("team_name", team.Name),
		zap.Int64("team_id", team.ID),
	)

	return nil
}

func (p *Payments) handleSubscriptionCreated(ctx context.Context, subscription stripe.Subscription) error {
	customerId := subscription.Customer.ID
	p.logCustomer(customerId, "Customer Subscription Created")
	return nil
}

func (p *Payments) handleSubscriptionUpdated(ctx context.Context, subscription stripe.Subscription) error {
	customerId := subscription.Customer.ID
	p.logCustomer(customerId, "Customer Subscription Updated")
	return nil
}

func (p *Payments) handleInvoicePaid(ctx context.Context, invoice stripe.Invoice) error {
	customerId := invoice.Customer.ID
	p.logCustomer(customerId, "Subscription Invoice Paid")

	team, err := p.database.GetTeamByStripeCustomerId(ctx, sql.NullString{Valid: true, String: customerId})
	if err != nil {
		return fmt.Errorf("Could not fetch team from DB using stripe customer ID %s: %s", customerId, err.Error())
	}

	subscriptionId := invoice.Subscription.ID

	if err != nil {
		return fmt.Errorf("Could not find subscription id %s: %s", subscriptionId, err.Error())
	}

	prod, err := p.GetSubscriptionProduct(subscriptionId)
	if err != nil {
		return fmt.Errorf("Could not fetch subscription %s products: %s", subscriptionId, err.Error())
	}

	credits, ok := prod.Metadata["credits_included"]
	if !ok {
		return fmt.Errorf("No credits included field in the product %s %s", prod.ID, prod.Name)
	}

	value, err := strconv.Atoi(credits)
	if err != nil {
		return fmt.Errorf(
			"Unable to parse included credits string in product (%s %s) '%s': %s",
			prod.ID, prod.Name, credits, err.Error())
	}

	if value <= 0 {
		return fmt.Errorf("Invalid amount of credits included (%d) in the subscription %s", value, subscriptionId)
	}

	sub_plan, err := p.database.AddSubscriptionCreditsByTeamId(ctx, database.AddSubscriptionCreditsByTeamIdParams{
		TeamID:           team.ID,
		RemainingCredits: int64(value),
	})

	if err != nil {
		return fmt.Errorf("Unable to add %d credits to team %d: %s", value, team.ID, err.Error())
	}

	sub_plan, err = p.database.SetSubscriptionStripeIdByTeamId(ctx, database.SetSubscriptionStripeIdByTeamIdParams{
		TeamID:               team.ID,
		StripeSubscriptionID: sql.NullString{Valid: true, String: invoice.Subscription.ID},
	})

	if err != nil {
		return fmt.Errorf("Unable to update stripe subscription id for team: %d: %s", team.ID, err.Error())
	}

	p.logger.Info(
		"Successfully granted credits to team",
		zap.Int("credits_added", value),
		zap.Int64("new_credit_amount", sub_plan.RemainingCredits),
		zap.String("team_name", team.Name),
		zap.Int64("team_id", team.ID),
	)

	return nil
}

func (p *Payments) handleInvoicePaymentFailed(ctx context.Context, invoice stripe.Invoice) error {
	customerId := invoice.Customer.ID
	p.logCustomer(customerId, "Subscription Invoice Payment Failed")
	return nil
}

func (p *Payments) logCustomer(stripeCustomerId string, event string) {

	customer, err := p.getCustomer(stripeCustomerId)

	if err != nil {
		p.logger.Error(
			"Could not fetch stripe customer",
			zap.String("stripe_customer_id", stripeCustomerId),
			zap.Error(err),
		)
		return
	}

	p.logger.Info(
		event,
		zap.String("customer_id", customer.ID),
		zap.String("customer_name", customer.Name),
		zap.String("customer_email", customer.Email),
	)
}

// RespondWithError sends an error response with a given HTTP status code and error message.
func (p *Payments) RespondWithError(w http.ResponseWriter, code int, message string) {
	p.logger.Error("Error occurred while handling stripe webhook", zap.String("error", message))
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
