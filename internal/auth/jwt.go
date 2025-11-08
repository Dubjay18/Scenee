package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKeyUserID struct{}

type JWTVerifier struct {
	Secret string
}

func NewJWTVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{Secret: secret}
}

func (v *JWTVerifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tok string

		// Try Authorization header first
		authz := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			tok = strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		} else {
			// Try cookie as fallback for browser requests
			if cookie, err := r.Cookie("access_token"); err == nil {
				tok = cookie.Value
			}
		}

		if tok == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parsed, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(v.Secret), nil
		})
		if err != nil || !parsed.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				r = r.WithContext(context.WithValue(r.Context(), ctxKeyUserID{}, sub))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func UserID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUserID{}).(string); ok {
		return v
	}
	return ""
}
