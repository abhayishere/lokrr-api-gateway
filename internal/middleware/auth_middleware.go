package middleware

import (
	"context"
	"net/http"

	authpb "github.com/abhayishere/lokrr-proto/gen/authpb"
)

const UserIDKey = "userID"

func AuthMiddleware(authClient authpb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			// todo: need to change token creation login and validate response according to it.
			res, err := authClient.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
				Token: token,
			})
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, res.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
