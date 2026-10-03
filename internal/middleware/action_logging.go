package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Action struct {
	UserId  string `json:"userId"`
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

type ActionLogging struct {
	ActionCache *ActionCache
}

func (a *ActionLogging) ActionLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload string

		if r.Body != nil && r.Body != http.NoBody {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			payload = string(bodyBytes)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		userId, _ := r.Context().Value(UserIdKey).(string)
		action := OutboxAction{
			UserId:    userId,
			Payload:   payload,
			Method:    r.Method, // Забираем метод (GET, POST и т.д.)
			CreatedAt: time.Now(),
		}
		a.ActionCache.Add(action)
		ctx := context.WithValue(r.Context(), UserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
