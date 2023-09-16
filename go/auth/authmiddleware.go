package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

var userCtxKey = &contextKey{"userId"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			clientKey := os.Getenv("CLERK_SECRET_KEY")

			if clientKey == "" {
				log.Fatalln("ERROR: CANNOT FIND CLERK CLIENT KEY")
			}

			client, _ := clerk.NewClient(clientKey)
			header := r.Header.Get("Authorization")

			// User is unauthenticated
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			sessionToken := strings.Split(header, " ")[1]
			sessClaims, err := client.VerifyToken(sessionToken)

			if err != nil {
				http.Error(w, "Invalid Authorization Token", http.StatusForbidden)
				return
			}

			user, err := client.Users().Read(sessClaims.Claims.Subject)

			if err != nil {
				http.Error(w, "Malformed Authorization Token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(ctx context.Context) *clerk.User {
	raw, _ := ctx.Value(userCtxKey).(*clerk.User)
	return raw
}

func EmailFromContext(ctx context.Context) (string, error) {
	user := FromContext(ctx)
	if user == nil {
		return "", fmt.Errorf("Not logged in")
	}
	for _, emailAddr := range user.EmailAddresses {
		if emailAddr.ID == *user.PrimaryEmailAddressID {
			return emailAddr.EmailAddress, nil
		}
	}
	return user.EmailAddresses[0].EmailAddress, nil
}

func FullnameFromContext(ctx context.Context) (string, error) {
	user := FromContext(ctx)
	if user == nil {
		return "", fmt.Errorf("Not logged in")
	}
	return fmt.Sprintf("%s %s", *user.FirstName, *user.LastName), nil
}
