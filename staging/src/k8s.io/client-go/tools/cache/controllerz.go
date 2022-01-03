package cache

import (
	"context"
)

type QueueKey interface {
	Namespace() string
	Name() string
}

type ControllerzConfig struct {
	ObjectKeyFunc KeyFunc
	DecodeKeyFunc func(key string) QueueKey

	ListAllIndex          string
	ListAllIndexFunc      IndexFunc
	ListAllIndexValueFunc func(ctx context.Context) (string, error)

	NamespaceIndex     string
	NamespaceIndexFunc IndexFunc
	NamespaceKeyFunc   func(ctx context.Context, namespace string) (string, error)

	NameKeyFunc          func(ctx context.Context, name string) (string, error)
	NamespaceNameKeyFunc func(ctx context.Context, namespace, name string) (string, error)

	NewSyncContextFunc func(ctx context.Context, key QueueKey) context.Context
}

type completedConfig struct {
	ControllerzConfig
}

var cc *completedConfig

func Complete(c ControllerzConfig) {
	cc = &completedConfig{
		ControllerzConfig: c,
	}
}

func ObjectKeyFunc(obj interface{}) (string, error) {
	return cc.ObjectKeyFunc(obj)
}

func DecodeKeyFunc(key string) QueueKey {
	return cc.DecodeKeyFunc(key)
}

func UseListAllIndex() bool {
	return cc.ListAllIndex != "" && cc.ListAllIndexFunc != nil && cc.ListAllIndexValueFunc != nil
}

func ListAllIndex() string {
	return cc.ListAllIndex
}

func ListAllIndexFunc() IndexFunc {
	return cc.ListAllIndexFunc
}

// func ListAllIndexValueFunc(ctx context.Context) (string, error) {
// 	return cc.ListAllIndexValueFunc(ctx)
// }

func NamespaceIndex2() string {
	return cc.NamespaceIndex
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

// var keyFunc KeyFunc = defaultKeyFunc
// var keyFuncSet = make(chan struct{})

// func Key(obj interface{}) (string, error) {
// 	return keyFunc(obj)
// }

// func defaultKeyFunc(obj interface{}) (string, error) {
// 	acc, err := meta.Accessor(obj)
// 	if err != nil {
// 		return "", err
// 	}
// 	return acc.GetNamespace() + "/" + acc.GetName(), nil
// }

// func DecodeKey(key string) QueueKey {
// 	return decodeKeyFunc(key)
// }

// type DecodeKeyFunc func(key string) QueueKey

// var decodeKeyFunc DecodeKeyFunc = defaultDecodeKeyFunc

// func defaultDecodeKeyFunc(key string) QueueKey {
// 	parts := strings.Split(key, "/")
// 	return &kubernetesQueueKey{
// 		namespace: parts[0],
// 		name:      parts[1],
// 	}
// }

type kubernetesQueueKey struct {
	namespace, name string
}

func (k *kubernetesQueueKey) Namespace() string {
	return k.namespace
}

func (k *kubernetesQueueKey) Name() string {
	return k.name
}

// func KeyFunc() KeyFunc {
// 	return keyFunc
// }

// func SetKeyFunc(f KeyFunc) {
// 	keyFunc = f
// 	close(keyFuncSet)
// }

// type AppendFunc func(interface{})

// var LH *ListerHelper

// func Setup(
// 	keyFunc KeyFunc,
// 	dkf DecodeKeyFunc,
// 	listAllIndex string,
// 	listAllIndexValueFunc func(ctx context.Context) string,
// 	listByNamespaceIndex string,
// 	namespaceKeyFunc func(ctx context.Context, ns string) string,
// ) {
// 	SetKeyFunc(keyFunc)
// 	decodeKeyFunc = dkf

// 	LH = NewListerHelper(
// 		listAllIndex,
// 		listAllIndexValueFunc,
// 		listByNamespaceIndex,
// 		namespaceKeyFunc,
// 	)
// }

// type ListerHelper struct {
// 	listAllIndex          string
// 	listAllIndexValueFunc func(ctx context.Context) string
// 	listByNamespaceIndex  string
// 	namespaceKeyFunc      func(ctx context.Context, ns string) string
// }

// func NewListerHelper(
// 	listAllIndex string,
// 	listAllIndexValueFunc func(ctx context.Context) string,
// 	listByNamespaceIndex string,
// 	namespaceKeyFunc func(ctx context.Context, ns string) string,
// ) *ListerHelper {
// 	return &ListerHelper{
// 		listAllIndex:          listAllIndex,
// 		listAllIndexValueFunc: listAllIndexValueFunc,
// 		listByNamespaceIndex:  listByNamespaceIndex,
// 		namespaceKeyFunc:      namespaceKeyFunc,
// 	}
// }

// func (lh *ListerHelper) ListAll(ctx context.Context, indexer Indexer, selector labels.Selector, appendFn AppendFunc) error {
// 	selectAll := selector.Empty()

// 	var (
// 		items []interface{}
// 		err   error
// 	)

// 	if lh.listAllIndex != "" && lh.listAllIndexValueFunc != nil {
// 		items, err = indexer.ByIndex(lh.listAllIndex, lh.listAllIndexValueFunc(ctx))
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		items = indexer.List()
// 	}

// 	for _, m := range items {
// 		if selectAll {
// 			// Avoid computing labels of the objects to speed up common flows
// 			// of listing all objects.
// 			appendFn(m)
// 			continue
// 		}
// 		metadata, err := meta.Accessor(m)
// 		if err != nil {
// 			return err
// 		}
// 		if selector.Matches(labels.Set(metadata.GetLabels())) {
// 			appendFn(m)
// 		}
// 	}
// 	return nil
// }

// func (lh *ListerHelper) ListAllByNamespace(ctx context.Context, indexer Indexer, namespace string, selector labels.Selector, appendFn AppendFunc) error {
// 	if namespace == metav1.NamespaceAll {
// 		return lh.ListAll(ctx, indexer, selector, appendFn)
// 	}

// 	selectAll := selector.Empty()

// 	var (
// 		items []interface{}
// 		err   error
// 	)

// 	nsIndex := lh.listByNamespaceIndex
// 	if nsIndex == "" {
// 		nsIndex = NamespaceIndex
// 	}

// 	ns := namespace
// 	if lh.namespaceKeyFunc != nil {
// 		ns = lh.namespaceKeyFunc(ctx, namespace)
// 	}

// 	items, err = indexer.ByIndex(nsIndex, ns)
// 	if err != nil {
// 		return err
// 	}

// 	for _, m := range items {
// 		if selectAll {
// 			// Avoid computing labels of the objects to speed up common flows
// 			// of listing all objects.
// 			appendFn(m)
// 			continue
// 		}
// 		metadata, err := meta.Accessor(m)
// 		if err != nil {
// 			return err
// 		}
// 		if selector.Matches(labels.Set(metadata.GetLabels())) {
// 			appendFn(m)
// 		}
// 	}

// 	return nil
// }

// type syncContextKey int

// const queueKey syncContextKey = iota

// type NewSyncContextFunc func(ctx context.Context, key QueueKey) context.Context

// var newSyncContextFunc NewSyncContextFunc = defaultNewSyncContextFunc

// func defaultNewSyncContextFunc(ctx context.Context, key QueueKey) context.Context {
// 	return ctx
// }

// func NewSyncContext(ctx context.Context, key QueueKey) context.Context {
// 	return newSyncContextFunc(ctx, key)
// }

// func SetNewSyncContextFunc(f NewSyncContextFunc) {
// 	newSyncContextFunc = f
// }

////// KCP specific

// type contextKey int

// const clusterKey contextKey = iota

// func ContextWithCluster(ctx context.Context, cluster string) context.Context {
// 	return context.WithValue(ctx, clusterKey, cluster)
// }

// func ClusterFromContext(ctx context.Context) string {
// 	v := ctx.Value(clusterKey)
// 	if v == nil {
// 		return ""
// 	}
// 	return v.(string)
// }

// func KCPNewSyncContext(ctx context.Context, key QueueKey
// ) context.Context {
// 	if kcpKey, ok := key.(KCPQueueKeyI); ok {
// 		return ContextWithCluster(ctx, kcpKey.ClusterName())
// 	}
// 	return ctx
// }

// type KCPQueueKeyI interface {
// 	QueueKey

// 	ClusterName() string
// }

// type KCPQueueKey struct {
// 	clusterName, namespace, name string
// }

// func NewKCPQueueKey(clusterName, namespace, name string) *KCPQueueKey {
// 	return &KCPQueueKey{
// 		clusterName: clusterName,
// 		namespace:   namespace,
// 		name:        name,
// 	}
// }

// func (k *KCPQueueKey) ClusterName() string {
// 	return k.clusterName
// }

// func (k *KCPQueueKey) Namespace() string {
// 	return k.namespace
// }

// func (k *KCPQueueKey) Name() string {
// 	return k.name
// }
