package config

import "time"

const (
	// SocketBufferSize set buffer size for websocket connection
	SocketBufferSize = 1024

	// MessageBufferSize set buffer size for websocket message
	MessageBufferSize = 256

	// MaxMessageSize set max size allow for Websocket message
	MaxMessageSize = 512

	// PongWait set limit wait for receive messages from client
	PongWait = 3 * time.Second

	// PingPeriod set time to ping client
	PingPeriod = PongWait * 9 / 10

	// WriteWait set limit wait for writing to client
	WriteWait = 3 * time.Second

	// ReadTimeout set wait time limit util we got request body
	ReadTimeout = 5 * time.Second
	// WriteTimeout set wait time limit util we got response
	WriteTimeout = 5 * time.Second
)
