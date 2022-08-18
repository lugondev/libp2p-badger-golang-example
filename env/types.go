package env

import "crypto/ecdsa"

// ConfigRaft configuration for raft node
type ConfigRaft struct {
	NodeId    string `mapstructure:"node_id"`
	Port      int    `mapstructure:"port"`
	VolumeDir string `mapstructure:"volume_dir"`
}

// ConfigServer configuration for server
type ConfigServer struct {
	Port       int `mapstructure:"port"`
	PrivateKey *ecdsa.PrivateKey
}

// Config configuration
type Config struct {
	Server ConfigServer `mapstructure:"server"`
	Raft   ConfigRaft   `mapstructure:"raft"`
}
