package calcium

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ricebook.net/platform/core/network/calico"
	"gitlab.ricebook.net/platform/core/scheduler/simple"
	"gitlab.ricebook.net/platform/core/source/gitlab"
	"gitlab.ricebook.net/platform/core/store/mock"
	"gitlab.ricebook.net/platform/core/types"
)

func TestGetRandomNode(t *testing.T) {
	store := &mockstore.MockStore{}
	config := types.Config{}
	c := &calcium{store: store, config: config, scheduler: simplescheduler.New(), network: calico.New(), source: gitlab.New(config)}

	n1 := &types.Node{Name: "node1", Podname: "podname", Endpoint: "tcp://10.0.0.1:2376", CPU: types.CPUMap{"0": 10, "1": 10}, Available: true}
	n2 := &types.Node{Name: "node2", Podname: "podname", Endpoint: "tcp://10.0.0.2:2376", CPU: types.CPUMap{"0": 10, "1": 10}, Available: true}

	store.On("GetNodesByPod", "podname").Return([]*types.Node{n1, n2}, nil)
	store.On("GetNode", "podname", "node1").Return(n1, nil)
	store.On("GetNode", "podname", "node2").Return(n2, nil)

	node, err := getRandomNode(c, "podname")
	assert.Contains(t, []string{"node1", "node2"}, node.Name)
	assert.Nil(t, err)
}

func TestGetRandomNodeFail(t *testing.T) {
	store := &mockstore.MockStore{}
	config := types.Config{}
	c := &calcium{store: store, config: config, scheduler: simplescheduler.New(), network: calico.New(), source: gitlab.New(config)}

	n1 := &types.Node{Name: "node1", Podname: "podname", Endpoint: "tcp://10.0.0.1:2376", CPU: types.CPUMap{"0": 10, "1": 10}, Available: false}
	n2 := &types.Node{Name: "node2", Podname: "podname", Endpoint: "tcp://10.0.0.2:2376", CPU: types.CPUMap{"0": 10, "1": 10}, Available: false}

	store.On("GetNodesByPod", "podname").Return([]*types.Node{n1, n2}, nil)
	store.On("GetNode", "podname", "node1").Return(n1, nil)
	store.On("GetNode", "podname", "node2").Return(n2, nil)

	node, err := getRandomNode(c, "podname")
	assert.NotNil(t, err)
	assert.Nil(t, node)
}

// 后面的我实在不想写了
// 让我们相信接口都是正确的吧
