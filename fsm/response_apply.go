package fsm

import "github.com/hashicorp/raft"

// ApplyResponse response from Apply raft
type ApplyResponse struct {
	Error  error
	Data   interface{}
	NodeID raft.ServerID
}
