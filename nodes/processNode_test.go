package nodes

import (
	"fmt"
	"testing"
)

type Log struct {
}

func (log *Log) Infof(format string, args ...interface{}) {
	fmt.Printf("test1 \n")
}

func TestLeaderElection1(t *testing.T) {
	forever := make(chan int)
	errCh := make(chan error)
	electedCh := make(chan bool)

	go func() {
		err := NewProcessNode(1, []string{"127.0.0.1:2181"}, errCh, electedCh, nil).Run()
		if err != nil {
			fmt.Printf("the error is %s \n", err.Error())
		}
	}()

	<-forever
}

func TestLeaderElection2(t *testing.T) {
	forever := make(chan int)
	errCh := make(chan error)
	electedCh := make(chan bool)

	go func() {
		err := NewProcessNode(2, []string{"127.0.0.1:2181"}, errCh, electedCh, nil).Run()
		if err != nil {
			fmt.Printf("the error is %s \n", err.Error())
		}
	}()

	<-forever
}

func TestLeaderElection3(t *testing.T) {
	forever := make(chan int)
	errCh := make(chan error)
	electedCh := make(chan bool)

	go func() {
		err := NewProcessNode(3, []string{"127.0.0.1:2181"}, errCh, electedCh, &Log{}).Run()
		if err != nil {
			fmt.Printf("the error is %s \n", err.Error())
		}
	}()

	<-forever
}

func TestLeaderElection4(t *testing.T) {
	forever := make(chan int)
	errCh := make(chan error)
	electedCh := make(chan bool)

	go func() {
		err := NewProcessNode(4, []string{"127.0.0.1:2181"}, errCh, electedCh, &Log{}).Run()
		if err != nil {
			fmt.Printf("the error is %s \n", err.Error())
		}
	}()

	<-forever
}
