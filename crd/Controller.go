package crd

import (
	"github.com/illublank/go-common/config"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	CrdController
}

func NewController(cfg config.Config, informer cache.SharedIndexInformer) *Controller {

	informer.GetStore()
	return nil
}
