package informers_test

import (
  "testing"

  "github.com/illublank/iron/informers"
  v1 "github.com/illublank/iron/test/typed/rdep/v1"
)

func TestInformers(t *testing.T) {

  informers.NewInformers[v1.ChannelList](nil, "test", nil)
}
