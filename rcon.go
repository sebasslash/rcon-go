package rcon_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type RCONClient struct {
	conn   net.Conn
	Config *RCONConfig
}

type RCONConfig struct {
	Host           string
	Password       string
	Port           uint32
	MaxPayloadSize int
}

const (
	BadAuth = -1
)

type packetType int32

const (
	PacketResponse packetType = iota
	_
	PacketCommand
	PacketLogin
)

type RCONHeader struct {
	Size      int32
	RequestID int32
	Type      packetType
}

func (c *RCONClient) Open() error {
	if c.Config.Host == "" {
		return ErrNoHostSpecified
	}

	log.Printf("[INFO] Dial connection %s\n", c.Config.Host)
	addr := fmt.Sprintf("%s:%d", c.Config.Host, c.Config.Port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return ErrFailedToConnect
	}

	c.conn = conn

	log.Printf("[INFO] Connection successfully established\n")
	return c.authenticate()
}

func (c *RCONClient) Close() error {
	return c.conn.Close()
}

func (c *RCONClient) UseDefaultValues() {
	if c.Config.Port == 0 {
		c.Config.Port = 27015
	}

	if c.Config.MaxPayloadSize == 0 {
		c.Config.MaxPayloadSize = 4096
	}
}

func (c *RCONClient) packetize(t packetType, p []byte) ([]byte, error) {
	pad := [2]byte{}
	length := int32(len(p) + 10)
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, length)
	_ = binary.Write(&buffer, binary.LittleEndian, int32(0))
	_ = binary.Write(&buffer, binary.LittleEndian, t)
	_ = binary.Write(&buffer, binary.LittleEndian, p)
	_ = binary.Write(&buffer, binary.LittleEndian, pad)

	if buffer.Len() >= c.Config.MaxPayloadSize {
		return nil, ErrPayloadLimitExceeded
	}

	return buffer.Bytes(), nil
}

func (c *RCONClient) depacketize(r io.Reader) (*RCONHeader, []byte, error) {
	header := &RCONHeader{}
	err := binary.Read(r, binary.LittleEndian, header)
	if err != nil {
		return nil, nil, err
	}

	payload := make([]byte, header.Size-8)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, nil, err
	}

	if header.Type != PacketResponse && header.Type != PacketCommand {
		return nil, nil, errors.New("bad packet type")
	}

	return header, payload[:len(payload)-2], nil
}

func (c *RCONClient) sendPacket(t packetType, p []byte) (*RCONHeader, []byte, error) {
	packet, err := c.packetize(t, p)
	if err != nil {
		return nil, nil, err
	}

	_, err = c.conn.Write(packet)
	if err != nil {
		return nil, nil, err
	}

	return c.depacketize(c.conn)
}

func (c *RCONClient) SendCommand(command string) (string, error) {
	if len([]byte(command)) > c.Config.MaxPayloadSize {
		return "", ErrPayloadLimitExceeded
	}

	header, payload, err := c.sendPacket(PacketCommand, []byte(command))
	if err != nil {
		return "", err
	}

	if header.RequestID == BadAuth {
		return "", ErrBadAuth
	}

	return string(payload), nil
}

func (c *RCONClient) authenticate() error {
	header, _, err := c.sendPacket(PacketLogin, []byte(c.Config.Password))
	if err != nil {
		return err
	}

	if header.RequestID == BadAuth {
		return ErrBadAuth
	}

	return nil
}
