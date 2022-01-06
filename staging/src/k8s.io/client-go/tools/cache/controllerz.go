package cache

import (
	"context"
)

type QueueKey interface {
	Namespace() string
	Name() string
}

type queueKey struct {
	namespace, name string
}

func (k *queueKey) Namespace() string {
	return k.namespace
}

func (k *queueKey) Name() string {
	return k.name
}

func DecodeMetaNamespaceKey(key string) (QueueKey, error) {
	ns, name, err := SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, err
	}

	return &queueKey{
		namespace: ns,
		name:      name,
	}, nil
}

type ControllerzConfig struct {
	ObjectKeyFunc KeyFunc
	DecodeKeyFunc func(key string) (QueueKey, error)

	ListAllIndexFunc      IndexFunc
	ListAllIndexValueFunc func(ctx context.Context) (string, error)

	NamespaceIndexFunc IndexFunc
	NamespaceKeyFunc   func(ctx context.Context, namespace string) (string, error)

	NameKeyFunc          func(ctx context.Context, name string) (string, error)
	NamespaceNameKeyFunc func(ctx context.Context, namespace, name string) (string, error)

	NewSyncContextFunc func(ctx context.Context, key QueueKey) context.Context
}

type completedConfig struct {
	ControllerzConfig
}

var (
	cc    *completedConfig
	ccSet = make(chan struct{})
)

func init() {
	const defaultListAllIndexValue = ""
	cc = &completedConfig{
		ControllerzConfig: ControllerzConfig{
			ObjectKeyFunc: DeletionHandlingMetaNamespaceKeyFunc,
			DecodeKeyFunc: DecodeMetaNamespaceKey,
			ListAllIndexFunc: func(obj interface{}) ([]string, error) {
				// Give all objects the same index value
				return []string{defaultListAllIndexValue}, nil
			},
			ListAllIndexValueFunc: func(ctx context.Context) (string, error) {
				// Match the index value from ListAllIndexFunc
				return defaultListAllIndexValue, nil
			},
			NamespaceIndexFunc: MetaNamespaceIndexFunc,
			NamespaceKeyFunc: func(ctx context.Context, namespace string) (string, error) {
				return namespace, nil
			},
			NameKeyFunc: func(ctx context.Context, name string) (string, error) {
				return name, nil
			},
			NamespaceNameKeyFunc: func(ctx context.Context, namespace, name string) (string, error) {
				if len(namespace) > 0 {
					return namespace + "/" + name, nil
				}
				return name, nil
			},
			NewSyncContextFunc: func(ctx context.Context, key QueueKey) context.Context {
				return ctx
			},
		},
	}
}

func Complete(c ControllerzConfig) {
	// This will panic if called twice
	close(ccSet)

	cc = &completedConfig{
		ControllerzConfig: c,
	}
}

func ObjectKeyFunc(obj interface{}) (string, error) {
	return cc.ObjectKeyFunc(obj)
}

func DecodeKeyFunc(key string) (QueueKey, error) {
	return cc.DecodeKeyFunc(key)
}

func ListAllIndexFunc() IndexFunc {
	return cc.ListAllIndexFunc
}

func NamespaceIndexFunc() IndexFunc {
	return cc.NamespaceIndexFunc
}

func NameKeyFunc(ctx context.Context, name string) (string, error) {
	return cc.NameKeyFunc(ctx, name)
}

func NamespaceNameKeyFunc(ctx context.Context, namespace, name string) (string, error) {
	return cc.NamespaceNameKeyFunc(ctx, namespace, name)
}

func NewSyncContext(ctx context.Context, key QueueKey) context.Context {
	return cc.NewSyncContextFunc(ctx, key)
}
