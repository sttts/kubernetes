package storage

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type Scope interface {
	NoNamespaceKeyRootFunc(prefix string) string
	PostDecode(obj runtime.Object) error
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
