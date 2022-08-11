package crd

import (
  "github.com/illublank/go-common/config"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/client-go/tools/cache"
  "k8s.io/client-go/tools/record"
)

type Controller struct {
  CrdController
  recorder record.EventRecorder
  informer cache.SharedIndexInformer
  Store    cache.Store
}

func NewController(cfg config.Config, gvr schema.GroupVersionResource, informer cache.SharedIndexInformer) *Controller {

  return &Controller{
    informer: informer,
    Store:    informer.GetStore(),
  }
}

func (s *Controller) AddEventHandler(handlerFuncs cache.ResourceEventHandlerFuncs) *Controller {

  s.informer.AddEventHandler(handlerFuncs)

  return s
}
