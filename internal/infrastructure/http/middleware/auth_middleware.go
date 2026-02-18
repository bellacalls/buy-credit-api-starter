package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sample-provider/buy-credit-api/internal/infrastructure/auth"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/response"
)

type contextKey string

const (
	PartnerIDKey contextKey = "partnerId"
	ClientIDKey  contextKey = "clientId"
)

type AuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "MISSING_AUTH_TOKEN", "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must be Bearer token")
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), PartnerIDKey, claims.PartnerID)
		ctx = context.WithValue(ctx, ClientIDKey, claims.ClientID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPartnerID(ctx context.Context) string {
	if partnerID, ok := ctx.Value(PartnerIDKey).(string); ok {
		return partnerID
	}
	return ""
}

func GetClientID(ctx context.Context) string {
	if clientID, ok := ctx.Value(ClientIDKey).(string); ok {
		return clientID
	}
	return ""
}
