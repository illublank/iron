package crd_test

import (
	"testing"

	"github.com/illublank/go-common/config/mock"
	"github.com/illublank/go-common/typ/collection"
	"github.com/illublank/iron/crd"
)

type sType struct {
	DeployStore  crd.DeploymentStore
	ServiceStore crd.ServiceStore
}

func TestController(t *testing.T) {
	cfg := mock.NewMapConfig(collection.NewGoMap())
	ctrl := crd.NewController[sType](cfg, nil, nil)

	ctrl.MockDiObj()

	// ctrl.AddEventHandler(func(obj interface{}, diObj *sType) {
	//   fmt.Println(diObj.DeployStore)
	// }, func(oldObj, newObj interface{}, diObj *sType) {}, func(obj interface{}, diObj *sType) {})
}
