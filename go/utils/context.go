package utils

type contextKey struct {
	name string
}

var DatabaseKey = &contextKey{name: "DatabaseContextKey"}
var AuthKey = &contextKey{name: "AuthContextKey"}
