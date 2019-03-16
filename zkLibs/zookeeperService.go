package zkLibs

import (
	"fmt"
	"os"
	"time"
)

type ZooKeeperService struct {
	Store
}

func NewZooKeeperService(urls []string, timeout time.Duration) (zk *ZooKeeperService, err error) {
	s, err := New(urls, &Config{
		ConnectionTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	zkSvr := &ZooKeeperService{
		Store: s,
	}

	return zkSvr, nil
}

func (zkService *ZooKeeperService) CreateNode(node string, watch bool, ephemeral bool) string {
	createNodePath := ""

	isExist, err := zkService.Store.Exists(node, &ReadOptions{
		Consistent: !ephemeral,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "check node %s exist error %s \n", node, err.Error())
		return createNodePath
	}

	if !isExist {
		err = zkService.Store.Put(node, []byte{}, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create node error %s \n", err.Error())
			return createNodePath
		}

		createNodePath = Normalize(node)
	} else {
		createNodePath = Normalize(node)
	}

	return createNodePath
}

func (zkService *ZooKeeperService) WatchNode(node string, runner EventRunner, stopCh chan struct{}, watch bool) chan error {
	_, err := zkService.Store.Exists(node, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "check node watched error %s \n", err.Error())
		return nil
	}
	errCh := zkService.Store.WatchNodeEvent(node, stopCh, runner, nil)

	return errCh
}

func (zkService *ZooKeeperService) GetChildren(node string, watch bool) []string {
	childNodes := []string{}

	kvPairs, err := zkService.Store.List(node, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot list the children error %s \n", err.Error())
	}

	for _, kvPair := range kvPairs {
		childNodes = append(childNodes, kvPair.Key)
	}

	return childNodes
}
