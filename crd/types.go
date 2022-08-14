package crd

import (
  "reflect"

  "k8s.io/client-go/tools/cache"
)

type DeployAndServiceStore struct {
  DeployStore  DeploymentStore
  ServiceStore ServiceStore
}

type DeploymentStore cache.Store

var DeploymentStoreType = reflect.TypeOf((*DeploymentStore)(nil)).Elem()

type PodStore cache.Store

var PodStoreType = reflect.TypeOf((*PodStore)(nil)).Elem()

type ServiceStore cache.Store

var ServiceStoreType = reflect.TypeOf((*ServiceStore)(nil)).Elem()

type MockStore struct {
  cache.Store
}

func (s *MockStore) Add(obj interface{}) error {
  return nil
}

func (s *MockStore) Update(obj interface{}) error {
  return nil
}

func (s *MockStore) Delete(obj interface{}) error {
  return nil
}

func (s *MockStore) List() []interface{} {
  return nil
}

func (s *MockStore) ListKeys() []string {
  return nil
}

func (s *MockStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
  return nil, false, nil
}

func (s *MockStore) GetByKey(key string) (item interface{}, exists bool, err error) {
  return nil, false, nil
}

func (s *MockStore) Replace([]interface{}, string) error {
  return nil
}

func (s *MockStore) Resync() error {
  return nil
}
