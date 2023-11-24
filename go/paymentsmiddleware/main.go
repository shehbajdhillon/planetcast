package paymentsmiddleware

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

type Payments struct {
	secretKey string
}

func Connect() *Payments {
	STRIPE_KEY := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = STRIPE_KEY
	return &Payments{secretKey: STRIPE_KEY}
}

// Customer Management
func (p *Payments) CreateCustomer(email string, paymentMethod string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email:         stripe.String(email),
		PaymentMethod: stripe.String(paymentMethod),
	}
	return customer.New(params)
}

func (p *Payments) GetCustomer(customerId string) (*stripe.Customer, error) {
	return customer.Get(customerId, nil)
}

func (p *Payments) UpdateCustomer(customerId string, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	return customer.Update(customerId, params)
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
