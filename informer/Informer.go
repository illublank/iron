package informer

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func NewInformers[ListT runtime.Object](client rest.Interface, namespace string, ctx context.Context, optsfuncs ...internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	var indexers cache.Indexers = cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}

	objListType := reflect.TypeOf((*ListT)(nil)).Elem().Elem()
	f, _ := objListType.FieldByName("Items")
	objType := f.Type.Elem()
	var obj runtime.Object = reflect.New(objType).Interface().(runtime.Object)
	resource := strings.ToLower(objType.Name()) + "s"
	fmt.Println(objListType, objType, resource)
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				for i := 0; i < len(optsfuncs); i++ {
					optsfuncs[i](&opts)
				}
				var result ListT = reflect.New(objListType).Interface().(ListT)
				var timeout time.Duration
				if opts.TimeoutSeconds != nil {
					timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
				}
				err := client.Get().Namespace(namespace).Resource(resource).VersionedParams(&opts, scheme.ParameterCodec).Timeout(timeout).Do(ctx).Into(result)
				return result, err
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				for i := 0; i < len(optsfuncs); i++ {
					optsfuncs[i](&opts)
				}
				var timeout time.Duration
				if opts.TimeoutSeconds != nil {
					timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
				}
				opts.Watch = true
				return client.Get().Namespace(namespace).Resource(resource).VersionedParams(&opts, scheme.ParameterCodec).Timeout(timeout).Watch(ctx)
			},
		},
		obj,
		30*time.Second,
		indexers,
	)
}
