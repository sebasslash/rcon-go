package rcon_go

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRCONClient_Open(t *testing.T) {
	client := createTestClient(t)
	err := client.Open()
	assert.NoError(t, err)
	defer client.Close()

	t.Run("when server connection refused", func(t *testing.T) {
		client2 := createTestClientInvalidHost(t)
		err := client2.Open()
		assert.Equal(t, ErrFailedToConnect, err)
	})

	t.Run("when bad auth", func(t *testing.T) {
		client3 := createTestClientNoAuth(t)
		err := client3.Open()
		assert.Equal(t, ErrBadAuth, err)
	})
}

func TestRCONClient_Close(t *testing.T) {
	client := createTestClient(t)

	t.Run("when client is open", func(t *testing.T) {
		assert.NoError(t, client.Open())
		assert.NoError(t, client.Close())
	})

	t.Run("when client is closed", func(t *testing.T) {
		err := client.Close()
		assert.NotNil(t, err)
	})
}

func TestRCONClient_SendCommand(t *testing.T) {
	// Because commands are unique depending on the server
	// we'll need to dynamically set one for our tests
	cmd := os.Getenv("RCON_GO_TEST_COMMAND")
	if cmd == "" {
		t.Skip("Skipping SendCommand test, no test command found")
	}

	client := createTestClient(t)

	t.Run("using test command", func(t *testing.T) {
		err := client.Open()
		assert.NoError(t, err)

		defer client.Close()

		resp, err := client.SendCommand(cmd)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
