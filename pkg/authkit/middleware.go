package authkit

import (
	"context"
	"net/http"
	"strings"

	"github.com/didar-platform/pkg/platform"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

// RequireAuth parses the Authorization: Bearer <token> header and attaches
// the decoded claims to the request context, or writes 401.
func RequireAuth(issuer *TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				platform.WriteError(w, platform.ErrUnauthorized)
				return
			}
			claims, err := issuer.Parse(token)
			if err != nil {
				platform.WriteError(w, platform.ErrUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth behaves like RequireAuth but continues the chain (without
// claims in context) when no/invalid token is present.
func OptionalAuth(issuer *TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if ok && token != "" {
				if claims, err := issuer.Parse(token); err == nil {
					r = r.WithContext(context.WithValue(r.Context(), claimsContextKey, claims))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireHost ensures the authenticated caller has host privileges.
// Must be chained after RequireAuth.
func RequireHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := FromContext(r.Context())
		if !ok || !claims.IsHost {
			platform.WriteError(w, platform.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}
