package raft

import (
	"math/rand"
	"time"
)

// raft election time (milliseconds)
const (
	minElectionTimeout = 150
	maxElectionTimeout = 300

	minBroadcastTimeout = 1
	maxBroadcastTimeout = 20
)

type raftTimer struct {
	t *time.Timer
}

func newRaftTimer() raftTimer {
	return raftTimer{
		t: time.NewTimer(randElectionMillisecond()),
	}
}

func (rt raftTimer) reset() bool {
	return rt.t.Reset(randElectionMillisecond())
}

func (rt raftTimer) stop() bool {
	return rt.t.Stop()
}

func (rt raftTimer) C() <-chan time.Time {
	return rt.t.C
}

func randElectionMillisecond() time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxElectionTimeout-minElectionTimeout) + minElectionTimeout
	return time.Duration(n) * time.Millisecond
}

func randBroadcastMillisecond() time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxBroadcastTimeout-minBroadcastTimeout) + minBroadcastTimeout
	return time.Duration(n) * time.Millisecond
}
