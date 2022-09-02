package informer_test

import (
	"testing"

	"github.com/illublank/iron/informer"
	v1 "github.com/illublank/iron/test/typed/rdep/v1"
)

func TestInformers(t *testing.T) {
	informer.NewInformers[*v1.ChannelList](nil, "test", nil)
}
