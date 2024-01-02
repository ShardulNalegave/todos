package auth

import "github.com/ShardulNalegave/todos/go/database"

type AuthState struct {
	IsAuth    bool
	SessionID string
	Session   database.Session
}
