# rcon-go

A simple RCON client written in Go 

**Note:** This client doesn't yet implement the full RCON protocol spec.

## Installation

```sh
go get -u github.com/sebasslash/rcon-go
```

## Example Usage

```go
import (
    "log"

    rcon "github.com/sebasslash/rcon-go"
)

func main() {
	config := &rcon.RCONConfig{
		Host: "my-awesome-server-ip",
		Port: 27015,
		Password: "my-awesome-server-pwd",
		MaxPayloadSize: 4096,
	}

    client := &rcon.RCONClient{
        Config: config
    }

    err := client.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Using a Mordhau specific RCON command 
    resp, err := client.SendCommand("playerlist")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Server response:\n %s", resp)
}
```

## Running Tests

In order to run the tests, a few env variables are expected:

1. `RCON_GO_HOST` - The host to test against 
1. `RCON_GO_PORT` - The port the RCON service is listening on 
1. `RCON_GO_PWD`  - The password used for authentication
1. `RCON_GO_TEST_COMMAND` - A test command to send to the server
