package calcium

// All functions are just proxy to store, since I don't want store to be exported.
// All these functions are meta data related.

import "gitlab.ricebook.net/platform/core/types"

func (c *calcium) ListPods() ([]*types.Pod, error) {
	return c.store.GetAllPods()
}

func (c *calcium) AddPod(podname, desc string) (*types.Pod, error) {
	return c.store.AddPod(podname, desc)
}

func (c *calcium) GetPod(podname string) (*types.Pod, error) {
	return c.store.GetPod(podname)
}

func (c *calcium) AddNode(nodename, endpoint, podname, cafile, certfile, keyfile string, public bool) (*types.Node, error) {
	return c.store.AddNode(nodename, endpoint, podname, cafile, certfile, keyfile, public)
}

func (c *calcium) GetNode(podname, nodename string) (*types.Node, error) {
	return c.store.GetNode(podname, nodename)
}

func (c *calcium) SetNodeAvailable(podname, nodename string, available bool) (*types.Node, error) {
	n, err := c.store.GetNode(podname, nodename)
	if err != nil {
		return nil, err
	}
	n.Available = available
	if err := c.store.UpdateNode(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (c *calcium) RemoveNode(nodename, podname string) (*types.Pod, error) {
	n, err := c.store.GetNode(podname, nodename)
	if err != nil {
		return nil, err
	}
	c.store.DeleteNode(n)
	return c.store.GetPod(podname)
}

func (c *calcium) ListPodNodes(podname string, all bool) ([]*types.Node, error) {
	var nodes []*types.Node
	candidates, err := c.store.GetNodesByPod(podname)
	if err != nil {
		return nodes, err
	}
	for _, candidate := range candidates {
		if candidate.Available || all {
			nodes = append(nodes, candidate)
		}
	}
	return nodes, err
}

func (c *calcium) GetContainer(id string) (*types.Container, error) {
	return c.store.GetContainer(id)
}

func (c *calcium) GetContainers(ids []string) ([]*types.Container, error) {
	return c.store.GetContainers(ids)
}
