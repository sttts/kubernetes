package rest

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scope interface {
	Name() string
	CacheKey(in string) string
	ScopeRequest(req *http.Request) error
}

type Scoper interface {
	// ScopeFromContext(ctx context.Context) (Scope, error)
	ScopeFromObject(obj metav1.Object) (Scope, error)
	ScopeFromKey(key string) (Scope, error)
}

type scopeKeyType int

const scopeKey scopeKeyType = iota

func WithScope(parent context.Context, scope Scope) context.Context {
	return context.WithValue(parent, scopeKey, scope)
}

func ScopeFrom(ctx context.Context) Scope {
	if v, ok := ctx.Value(scopeKey).(Scope); ok {
		return v
	}
	return nil
}
