package crd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/illublank/go-common/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Controller[T any] struct {
  CrdController
  informer        cache.SharedIndexInformer
  informerFactory informers.SharedInformerFactory
  client          kubernetes.Interface
  diObj           *T
  stopCh          chan struct{}
}

func NewController[T any](cfg config.Config, client kubernetes.Interface, gvr schema.GroupVersionResource, informer cache.SharedIndexInformer) *Controller[T] {

  return &Controller[T]{
    client:   client,
    informer: informer,
    stopCh:   make(chan struct{}),
  }
}

func (s *Controller[T]) EnableK8sResource(sync time.Duration, options ...informers.SharedInformerOption) *ControllerWithDiObj[T] {
  s.informerFactory = informers.NewSharedInformerFactoryWithOptions(s.client, sync, options...)
  typ := reflect.TypeOf(s.diObj).Elem()
  valPtr := reflect.New(typ)
  val := valPtr.Elem()
  for i := 0; i < typ.NumField(); i++ {
    f := typ.Field(i)
    t := f.Type
    switch {
    case t.AssignableTo(DeploymentStoreType):
      val.Field(i).Set(reflect.ValueOf(DeploymentStore(s.informerFactory.Extensions().V1beta1().Deployments().Informer().GetStore())))
    case t.AssignableTo(ServiceStoreType):
      val.Field(i).Set(reflect.ValueOf(ServiceStore(s.informerFactory.Core().V1().Services().Informer().GetStore())))
    case t.AssignableTo(PodStoreType):
      val.Field(i).Set(reflect.ValueOf(PodStore(s.informerFactory.Core().V1().Pods().Informer().GetStore())))
    default:
    }
  }
  s.diObj = valPtr.Interface().(*T)
  return &ControllerWithDiObj[T]{
    c: s,
    diObj: s.diObj,
  }
}

func (s *Controller[T]) MockDiObj() {
  typ := reflect.TypeOf(s.diObj).Elem()
  valPtr := reflect.New(typ)
  val := valPtr.Elem()
  fmt.Println("diObj type:", typ)
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
  s.diObj = valPtr.Interface().(*T)
}

func (s *Controller[T]) AddEventHandler(
  addFunc func(obj interface{}),
  updateFunc func(oldObj, newObj interface{}),
  deleteFunc func(obj interface{}),
) *Controller[T] {

  s.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: addFunc,
    UpdateFunc: updateFunc,
    DeleteFunc: deleteFunc,
  })

  return s
}

func (s *Controller[T]) Run() {
  s.informerFactory.Start(s.stopCh)
  s.informerFactory.WaitForCacheSync(s.stopCh)
}

func (s *Controller[T]) Stop() {
  s.stopCh <- struct{}{}
}

type ControllerWithDiObj[T any] struct {
  c *Controller[T]
  diObj *T
}

func (s *ControllerWithDiObj[T]) AddEventHandler(
  addFunc func(obj interface{}, diObj *T),
  updateFunc func(oldObj, newObj interface{}, diObj *T),
  deleteFunc func(obj interface{}, diObj *T),
) *Controller[T] {

  s.c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
      addFunc(obj, s.diObj)
    },
    UpdateFunc: func(oldObj, newObj interface{}) {
      updateFunc(oldObj, newObj, s.diObj)
    },
    DeleteFunc: func(obj interface{}) {
      deleteFunc(obj, s.diObj)
    },
  })

  return s.c
}