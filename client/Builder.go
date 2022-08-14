package client

import (
  "os"

  "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
  "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"

  // "k8s.io/client-go/kubernetes/scheme"
  "k8s.io/client-go/rest"
)

type K8sClientBuilder interface {
  GetConfig() *rest.Config
  GetClientset() (kubernetes.Interface, error)
  GetExtClientset() (clientset.Interface, error)
  GetRestClient(*schema.GroupVersion, bool) (rest.Interface, error)
}

type K8sClientBuilderImpl struct {
  K8sClientBuilder
  config *rest.Config
}

func (s *K8sClientBuilderImpl) GetConfig() *rest.Config {
  return s.config
}

func (s *K8sClientBuilderImpl) GetClientset() (kubernetes.Interface, error) {
  return kubernetes.NewForConfig(s.config)
}

func (s *K8sClientBuilderImpl) GetExtClientset() (clientset.Interface, error) {
  return clientset.NewForConfig(s.config)
}

func (s *K8sClientBuilderImpl) GetRestClient(gv *schema.GroupVersion, unversion bool) (rest.Interface, error) {
  cfg := s.GetConfig()
  cfg.ContentConfig.GroupVersion = gv
  cfg.APIPath = "/apis"
  cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
  cfg.UserAgent = rest.DefaultKubernetesUserAgent()
  if unversion {
    return rest.UnversionedRESTClientFor(cfg)
  }
  return rest.RESTClientFor(cfg)
}

func NewK8sClientBuilder() (K8sClientBuilder, error) {

  cfg, err := func() (*rest.Config, error) {
    kubeConfig := os.Getenv("KUBECONFIG")
    if len(kubeConfig) != 0 {
      return clientcmd.BuildConfigFromFlags("", kubeConfig)
    } else {
      return rest.InClusterConfig()
    }
  }()

  if err != nil {
    return nil, err
  }
  return &K8sClientBuilderImpl{
    config: cfg,
  }, nil
}
