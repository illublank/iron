package crd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/illublank/go-common/config"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Controller[T any] struct {
	CrdController
	informer              cache.SharedIndexInformer
	informerStopCh        chan struct{}
	informerFactory       informers.SharedInformerFactory
	informerFactoryStopCh chan struct{}
	client                kubernetes.Interface
	stores                *T
	stopCh                chan struct{}
	doneCh                chan struct{}
}

func NewController[T any](cfg config.Config, restClient rest.Interface, informer cache.SharedIndexInformer) *Controller[T] {
	return &Controller[T]{
		client:         kubernetes.New(restClient),
		informer:       informer,
		informerStopCh: make(chan struct{}),
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}
}

func (s *Controller[T]) InjectStores(sync time.Duration, options ...informers.SharedInformerOption) *ControllerWithStores[T] {
	s.informerFactory = informers.NewSharedInformerFactoryWithOptions(s.client, sync, options...)
	s.informerFactoryStopCh = make(chan struct{})
	typ := reflect.TypeOf(s.stores).Elem()
	valPtr := reflect.New(typ)
	val := valPtr.Elem()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		t := f.Type
		switch {
		case t.AssignableTo(DeploymentStoreType):
			val.Field(i).Set(reflect.ValueOf(DeploymentStore(s.informerFactory.Apps().V1().Deployments().Informer().GetStore())))
		case t.AssignableTo(ServiceStoreType):
			val.Field(i).Set(reflect.ValueOf(ServiceStore(s.informerFactory.Core().V1().Services().Informer().GetStore())))
		case t.AssignableTo(PodStoreType):
			val.Field(i).Set(reflect.ValueOf(PodStore(s.informerFactory.Core().V1().Pods().Informer().GetStore())))
		case t.AssignableTo(ConfigMapStoreType):
			val.Field(i).Set(reflect.ValueOf(ConfigMapStore(s.informerFactory.Core().V1().ConfigMaps().Informer().GetStore())))
		default:
		}
	}
	s.stores = valPtr.Interface().(*T)
	return &ControllerWithStores[T]{
		c:      s,
		stores: s.stores,
	}
}

func (s *Controller[T]) MockDiObj() {
	typ := reflect.TypeOf(s.stores).Elem()
	valPtr := reflect.New(typ)
	val := valPtr.Elem()
	fmt.Println("stores type:", typ)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		t := f.Type
		switch {
		case t.AssignableTo(DeploymentStoreType):
			fmt.Println("field:", f)
			val.Field(i).Set(reflect.ValueOf(DeploymentStore(&MockStore{})))
		case t.AssignableTo(ServiceStoreType):
			fmt.Println("field:", f)
			val.Field(i).Set(reflect.ValueOf(ServiceStore(&MockStore{})))
		case t.AssignableTo(PodStoreType):
			fmt.Println("field:", f)
			val.Field(i).Set(reflect.ValueOf(PodStore(&MockStore{})))
		default:
		}
	}
	s.stores = valPtr.Interface().(*T)
}

func (s *Controller[T]) AddEventHandler(
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{}),
) *Controller[T] {
	s.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	})

	return s
}

func (s *Controller[T]) Run() {
	go func() {
		<-s.stopCh
		s.informerStopCh <- struct{}{}
		if s.informerFactoryStopCh != nil {
			s.informerFactoryStopCh <- struct{}{}
		}
	}()
	if s.informerFactory != nil {
		s.informerFactory.Start(s.informerFactoryStopCh)
		s.informerFactory.WaitForCacheSync(s.informerFactoryStopCh)
	}
	s.informer.Run(s.informerStopCh)
	s.doneCh <- struct{}{}
}

func (s *Controller[T]) Stop() {
	s.stopCh <- struct{}{}
	<-s.doneCh
}

type ControllerWithStores[T any] struct {
	c      *Controller[T]
	stores *T
}

func (s *ControllerWithStores[T]) AddEventHandler(
	addFunc func(obj interface{}, stores *T),
	updateFunc func(oldObj, newObj interface{}, stores *T),
	deleteFunc func(obj interface{}, stores *T),
) *Controller[T] {
	s.c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addFunc(obj, s.stores)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateFunc(oldObj, newObj, s.stores)
		},
		DeleteFunc: func(obj interface{}) {
			deleteFunc(obj, s.stores)
		},
	})

	return s.c
}
