package raft_handler

import (
	"github.com/hashicorp/raft"
)

// handler struct handler
type handlerRaft struct {
	raft *raft.Raft
}

func New(raft *raft.Raft) *handlerRaft {
	return &handlerRaft{
		raft: raft,
	}
}
