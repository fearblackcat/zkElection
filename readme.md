# zookeeper

## Leader Election Encapsulation 

A simple way of doing leader election with ZooKeeper is to use the SEQUENCE|EPHEMERAL flags when creating znodes that represent "proposals" of clients. The idea is to have a znode, say "/election", such that each znode creates a child znode "/election/n_" with both flags SEQUENCE|EPHEMERAL. With the sequence flag, ZooKeeper automatically appends a sequence number that is greater that any one previously appended to a child of "/election". The process that created the znode with the smallest appended sequence number is the leader.

That's not all, though. It is important to watch for failures of the leader, so that a new client arises as the new leader in the case the current leader fails. A trivial solution is to have all application processes watching upon the current smallest znode, and checking if they are the new leader when the smallest znode goes away (note that the smallest znode will go away if the leader fails because the node is ephemeral). But this causes a herd effect: upon of failure of the current leader, all other processes receive a notification, and execute getChildren on "/election" to obtain the current list of children of "/election". If the number of clients is large, it causes a spike on the number of operations that ZooKeeper servers have to process. To avoid the herd effect, it is sufficient to watch for the next znode down on the sequence of znodes. If a client receives a notification that the znode it is watching is gone, then it becomes the new leader in the case that there is no smaller znode. Note that this avoids the herd effect by not having all clients watching the same znode.

Here's the pseudo code:

Let ELECTION be a path of choice of the application. To volunteer to be a leader:

    Create znode z with path "ELECTION/n_" with both SEQUENCE and EPHEMERAL flags;

    Let C be the children of "ELECTION", and i be the sequence number of z;

    Watch for changes on "ELECTION/n_j", where j is the largest sequence number such that j < i and n_j is a znode in C;

Upon receiving a notification of znode deletion:

    Let C be the new set of children of ELECTION;

    If z is the smallest node in C, then execute leader procedure;

    Otherwise, watch for changes on "ELECTION/n_j", where j is the largest sequence number such that j < i and n_j is a znode in C;

Note that the znode having no preceding znode on the list of children does not imply that the creator of this znode is aware that it is the current leader. Applications may consider creating a separate znode to acknowledge that the leader has executed the leader procedure. 

## Usage

### Reference the nodes test file

```shell
$ cd $GOPATH/src/github.com/fearblackcat/
$ git clone https://github.com/fearblackcat/zkElection.git
$ cd ./zkElection/nodes
$ go test -v -run ^TestLeaderElection1$
=== RUN   TestLeaderElection1
process id 1 has started
2019/02/27 14:47:53 Connected to 127.0.0.1:2181
2019/02/27 14:47:53 Authenticated: id=72057738101588070, timeout=4000
2019/02/27 14:47:53 Re-submitting `0` credentials after reconnect
process id is 1 and path is /election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075
the childNodePaths is [election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075] and ProcessNodePath is /election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075
process id is 1 and get the leader right

$ go test -v -run ^TestLeaderElection2$
=== RUN   TestLeaderElection2
process id 2 has started
2019/02/27 14:47:57 Connected to 127.0.0.1:2181
2019/02/27 14:47:57 Authenticated: id=72057738101588071, timeout=4000
2019/02/27 14:47:57 Re-submitting `0` credentials after reconnect
process id is 2 and path is /election/_c_ca49ed4caebf9d4c8b8838f283c4dbbd-p_0000000076
the childNodePaths is [election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075 election/_c_ca49ed4caebf9d4c8b8838f283c4dbbd-p_0000000076] and ProcessNodePath is /election/_c_ca49ed4caebf9d4c8b8838f283c4dbbd-p_0000000076
the watchedNodeShortPath is election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075
the process id 2 watch on node with path election/_c_2aabdd6e6c540cec7b8d68ea9113eaa3-p_0000000075


```
#### killed test1 and test2 will get the leader position
