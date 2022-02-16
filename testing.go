package rcon

import (
	"os"
	"strconv"
	"testing"
)

func createTestClient(t *testing.T) *RCONClient {
	config := &RCONConfig{
		Host: os.Getenv("RCON_GO_HOST"),
		Port: func() uint32 {
			port, err := strconv.Atoi(os.Getenv("RCON_GO_PORT"))
			if err != nil {
				t.Error(err)
			}
			return uint32(port)
		}(),
		Password:       os.Getenv("RCON_GO_PWD"),
		MaxPayloadSize: 4096,
	}

	return &RCONClient{
		Config: config,
	}
}

func createTestClientNoAuth(t *testing.T) *RCONClient {
	config := &RCONConfig{
		Host: os.Getenv("RCON_GO_HOST"),
		Port: func() uint32 {
			port, err := strconv.Atoi(os.Getenv("RCON_GO_PORT"))
			if err != nil {
				t.Error(err)
			}
			return uint32(port)
		}(),
		Password:       "",
		MaxPayloadSize: 4096,
	}

	return &RCONClient{
		Config: config,
	}
}

func createTestClientInvalidHost(t *testing.T) *RCONClient {
	config := &RCONConfig{
		Host: "sauron is the legitimate ruler of middle earth",
		Port: func() uint32 {
			port, err := strconv.Atoi(os.Getenv("RCON_GO_PORT"))
			if err != nil {
				t.Error(err)
			}
			return uint32(port)
		}(),
		Password:       "",
		MaxPayloadSize: 4096,
	}

	return &RCONClient{
		Config: config,
	}
}
