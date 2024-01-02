package auth

import (
	"context"
	"net/http"

	"github.com/ShardulNalegave/todos/go/database"
	"github.com/ShardulNalegave/todos/go/utils"
	"gorm.io/gorm"
)

const AuthCookie = "auth-session"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
			cookie, err := r.Cookie(AuthCookie)
			if err != nil {
				ctx := context.WithValue(r.Context(), utils.AuthKey, AuthState{IsAuth: false})
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			var session database.Session
			res := db.First(&session, "id = ?", cookie.Value)
			if res.Error != nil {
				ctx := context.WithValue(r.Context(), utils.AuthKey, AuthState{IsAuth: false})
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), utils.AuthKey, AuthState{
				IsAuth:    true,
				SessionID: session.ID,
				Session:   session,
			})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
