package nodes

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fearblackcat/zkElection/zkLibs"
)

type Logger interface {
	Infof(format string, args ...interface{})
}

type Runner interface {
	Run() error
}

type ProcessNode struct {
	Log             Logger
	ProcessNodePath string
	WatchNodePath   string
	ID              uint64
	ZkService       *zkLibs.ZooKeeperService
	StopCh          chan struct{}
	ErrCh           chan error
	ElectedCh       chan bool
}

const (
	LeaderRootNode    = "election"
	ProcessNodePrefix = "/p_"
)

func NewProcessNode(id uint64, zkURLS []string, errCh chan error, electedCh chan bool, dLog Logger) *ProcessNode {
	zookeeperService, err := zkLibs.NewZooKeeperService(zkURLS, 1*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot new zk service %s", err.Error())
	}

	node := &ProcessNode{
		Log:       dLog,
		ID:        id,
		ZkService: zookeeperService,
		StopCh:    make(chan struct{}),
		ErrCh:     errCh,
		ElectedCh: electedCh,
	}

	return node
}

func (node *ProcessNode) Run() error {
	if node.Log != nil {
		node.Log.Infof("process id %d has started", node.ID)
	} else {
		fmt.Printf("process id %d has started \n", node.ID)
	}

	if node.ZkService == nil {
		return errors.New("cannot initialize zookeeper service")
	}

	rootNodePath := node.ZkService.CreateNode(LeaderRootNode, false, false)
	if rootNodePath == "" {
		return errors.New("cannot create the leader root node")
	}

	nodePath, err := node.ZkService.Store.CreateSequenceEphemeralPath(LeaderRootNode+ProcessNodePrefix, []byte{})
	if err != nil {
		fmt.Fprint(os.Stderr, "create sequence ephemeral node failed")
		return err
	}

	node.ProcessNodePath = nodePath

	if node.Log != nil {
		node.Log.Infof("process id is %d and path is %s", node.ID, node.ProcessNodePath)
	} else {
		fmt.Printf("process id is %d and path is %s \n", node.ID, node.ProcessNodePath)
	}

	node.AttemptForLeaderPosition()

	return nil
}

func (node *ProcessNode) sendElectionSignal(isElection bool) {
	go func() {
		if node.ElectedCh != nil {
			node.ElectedCh <- isElection
		}
	}()
}

func (node *ProcessNode) AttemptForLeaderPosition() {
	childNodePaths := node.ZkService.GetChildren(LeaderRootNode, false)
	if len(childNodePaths) == 0 {
		fmt.Fprint(os.Stderr, "the leaderRootNode children is Empty")
		node.sendElectionSignal(false)
		return
	}

	sort.Sort(zkLibs.SequenceStrings(childNodePaths))

	if node.Log != nil {
		node.Log.Infof("the childNodePaths is %v and ProcessNodePath is %s", childNodePaths, node.ProcessNodePath)
	} else {
		fmt.Printf("the childNodePaths is %v and ProcessNodePath is %s \n", childNodePaths, node.ProcessNodePath)
	}

	if strings.Compare("/"+childNodePaths[0], node.ProcessNodePath) == 0 {
		if node.Log != nil {
			node.Log.Infof("process id is %d and get the leader right", node.ID)
		} else {
			fmt.Printf("process id is %d and get the leader right \n", node.ID)
		}
		node.sendElectionSignal(true)
	} else {
		index := -1
		node.sendElectionSignal(false)
		for i, path := range childNodePaths {
			if strings.Compare("/"+path, node.ProcessNodePath) == 0 {
				index = i
				break
			}
		}
		if index != -1 {
			watchedNodePath := childNodePaths[index-1]

			if node.Log != nil {
				node.Log.Infof("the watchedNodeShortPath is %s", watchedNodePath)
			} else {
				fmt.Printf("the watchedNodeShortPath is %s \n", watchedNodePath)
			}

			if node.Log != nil {
				node.Log.Infof("the process id %d watch on node with path %s", node.ID, watchedNodePath)
			} else {
				fmt.Printf("the process id %d watch on node with path %s \n", node.ID, watchedNodePath)
			}

			go func() {
				errCh := node.ZkService.WatchNode(watchedNodePath, node, node.StopCh, true)
				select {
				case err, ok := <-errCh:
					if ok && node.ErrCh != nil {
						node.ErrCh <- err
					}
				default:
				}
			}()
		}
	}

}

func (node *ProcessNode) Process(event *zkLibs.EventData) {
	if node.Log != nil {
		node.Log.Infof("process id is %d and event is %v", node.ID, *event)
	} else {
		fmt.Printf("process id is %d and event is %v \n", node.ID, *event)
	}

	if event.Event.Type.String() == "EventNodeDeleted" {
		node.AttemptForLeaderPosition()
	}
}
