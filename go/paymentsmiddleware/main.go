package paymentsmiddleware

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/graph/model"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
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

func (p *Payments) GetLast4CardDigits(paymentMethodId string) (string, error) {
	pm, err := paymentmethod.Get(paymentMethodId, nil)
	if err != nil {
		return "", nil
	}
	card := pm.Card
	return card.Last4, nil
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

func (p *Payments) DeleteCustomerFromDB(ctx context.Context, customerId string) error {
	team, err := p.database.GetTeamByStripeCustomerId(ctx, sql.NullString{Valid: true, String: customerId})

	if err != nil {
		return err
	}

	team, err = p.database.UpdateTeamStripeCustomerIdByTeamId(ctx, database.UpdateTeamStripeCustomerIdByTeamIdParams{
		ID:               team.ID,
		StripeCustomerID: sql.NullString{Valid: false, String: ""},
	})

	if err != nil {
		return err
	}

	p.logger.Info("Deleted stripe customer id for team", zap.String("team_name", team.Name), zap.Int64("team_id", team.ID))

	return nil
}

// Subscription Management
func (p *Payments) CreateSubscription(customerId string, priceId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerId),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(priceId)},
		},
	}
	return subscription.New(params)
}

func (p *Payments) GetSubscription(subscriptionId string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionId, nil)
}

func (p *Payments) GetSubscriptionProduct(subscriptionId string) (*stripe.Product, error) {

	sub, err := p.GetSubscription(subscriptionId)

	if err != nil {
		return nil, err
	}

	if len(sub.Items.Data) <= 0 {
		return nil, fmt.Errorf("No items found for subscription ID %s", subscriptionId)
	}

	firstItem := sub.Items.Data[0]
	productObject, err := product.Get(firstItem.Price.Product.ID, nil)
	if err != nil {
		return nil, err
	}

	return productObject, nil
}

func (p *Payments) GetSubscriptionAmountInUSDAndInterval(subscriptionId string) (int64, string, error) {
	sub, err := p.GetSubscription(subscriptionId)

	if err != nil {
		return 0, "", err
	}

	if len(sub.Items.Data) <= 0 {
		return 0, "", fmt.Errorf("No items found for subscription ID %s", subscriptionId)
	}

	firstItem := sub.Items.Data[0]
	interval := string(firstItem.Plan.Interval)

	// This assumes that firstItem.Plan.Currency will be USD
	// We are only listing out subscriptions in USD so this works for now
	// Will have to change in future when we list subsctipions in other currencies
	// currency := firstItem.Plan.Currency
	amount := firstItem.Plan.Amount

	return amount / 100, interval, nil
}

func (p *Payments) CancelSubscription(subscriptionId string) (*stripe.Subscription, error) {
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

func (p *Payments) GetSubscriptionInterval(subscriptionId string) (stripe.PlanInterval, error) {
	subPlan, err := p.GetSubscription(subscriptionId)

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", fmt.Errorf("No items found for subscription ID %s", subscriptionId)
	}

	firstItem := subPlan.Items.Data[0]
	interval := firstItem.Plan.Interval

	return interval, nil
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

func (p *Payments) GetSubscriptionPlanData(subscriptionId string) (*model.SubscriptionData, error) {

	subPlan, err := p.GetSubscription(subscriptionId)

	if err != nil {
		return nil, err
	}

	if len(subPlan.Items.Data) <= 0 {
		return nil, fmt.Errorf("No items found for subscription ID %s", subscriptionId)
	}

	firstItem := subPlan.Items.Data[0]
	currentPeriodStart := time.Unix(subPlan.CurrentPeriodStart, 0).UTC()
	currentPeriodEnd := time.Unix(subPlan.CurrentPeriodEnd, 0).UTC()
	status := string(subPlan.Status)
	interval := string(firstItem.Plan.Interval)

	// This assumes that firstItem.Plan.Currency will be USD
	// We are only listing out subscriptions in USD so this works for now
	// Will have to change in future when we list subsctipions in other currencies
	// currency := firstItem.Plan.Currency
	amount := firstItem.Plan.Amount
	costInUsd := (amount / 100)

	prod, err := p.GetSubscriptionProduct(subscriptionId)
	if err != nil {
		return nil, err
	}

	planName := prod.Name

	paymentMethodId := subPlan.DefaultPaymentMethod.ID
	lastFourCardDigits, err := p.GetLast4CardDigits(paymentMethodId)
	if err != nil {
		return nil, err
	}

	subscriptionData := model.SubscriptionData{
		CurrentPeriodStart: currentPeriodStart.String(),
		CurrentPeriodEnd:   currentPeriodEnd.String(),
		Status:             status,
		Interval:           interval,
		PlanName:           planName,
		CostInUsd:          costInUsd,
		LastFourCardDigits: lastFourCardDigits,
	}

	return &subscriptionData, nil
}
