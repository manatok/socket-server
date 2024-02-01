package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Producer struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Producer) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		dataPacket, err := unpackData(message)

		if err != nil {
			log.Printf("error: %v", err)
		} else {
			c.hub.broadcast <- dataPacket
		}
	}
}

func unpackData(data []byte) (dataPacket *DataPacket, err error) {

	fmt.Printf("Data: %s\n", data)
	c := DataPacket{}
	err = json.Unmarshal(data, &c)

	if err != nil {
		panic(fmt.Sprintf("Unmarshal error: %v\n", err))
	}

	fmt.Printf("Unmarshalled %+v\n", c)

	return &c, err
}
