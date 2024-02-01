package main

import (
	"github.com/golang/glog"
	"sync"
)

var hubMutex = sync.RWMutex{}

// Hub maintains a list of clients that are subscribed to various topics.
type Hub struct {
	// Registered clients. subscriptions[topicName][client]
	subscriptions map[string]map[*Client]bool

	// Channel used to register new clients on the hub.
	register chan *Client

	// Channel to un-register clients from the hub.
	unregister chan *Client

	// Inbound messages from the producers.
	broadcast chan *DataPacket
}

// JSON messages from the producer need to be in this format.
type DataPacket struct {
	Topic string `json:"Topic"`
	Data  string `json:"Data"`
}

// Factory method for creating a new Hub
func newHub() *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		subscriptions: make(map[string]map[*Client]bool),
		broadcast:     make(chan *DataPacket),
	}
}

// Run the main hub for client/message management.
func (h *Hub) run() {
	glog.Debug("main", STARTING_HUB)

	for {
		select {
		case client := <-h.register:
			registerClient(h, client)
		case client := <-h.unregister:
			unregisterClient(h, client)
		case message := <-h.broadcast:
			broadcast(h, message)
		}
	}
}

// Add a client to the hub
func registerClient(hub *Hub, client *Client) {
	glog.Debug("main", REGISTERING_CLIENT, client)

	hubMutex.Lock()
	defer hubMutex.Unlock()

	if _, ok := hub.subscriptions[client.topic]; !ok {
		hub.subscriptions[client.topic] = make(map[*Client]bool)
	}

	hub.subscriptions[client.topic][client] = true

	glog.Debug("main", TOTAL_TOPIC_CLIENTS, client.topic, len(hub.subscriptions[client.topic]))
}

// Remove a client from the hub
func unregisterClient(hub *Hub, client *Client) {
	glog.Debug("main", UNREGISTERING_CLIENT, client)

	hubMutex.Lock()
	defer hubMutex.Unlock()

	removeClient(hub, client)
}

// Send message to all clients subscribed to topic
func broadcast(hub *Hub, message *DataPacket) {
	hubMutex.RLock()
	defer hubMutex.RUnlock()

	for client := range hub.subscriptions[message.Topic] {
		select {
		case client.send <- []byte(message.Data):
			glog.Debug("main", MESSAGE_BROADCAST, message.Topic, len(message.Data))
		default:
			glog.Debug("main", CLIENT_BUFFER_FULL, client)

			removeClient(hub, client)
		}
	}
}

/**
 * NB: This should only be called if you currently have the hubMutex lock.
 *
 * This method has been created to factor out common code required by boradcast and unregisterClient.
 * If this code lived in unregisterClient and was called by broadcast a deadlock would occur.
 */
func removeClient(hub *Hub, client *Client) {
	if _, ok := hub.subscriptions[client.topic][client]; ok {
		delete(hub.subscriptions[client.topic], client)

		if len(hub.subscriptions[client.topic]) == 0 {
			glog.Debug("main", UNREGISTERING_TOPIC, client.topic)
			delete(hub.subscriptions, client.topic)
			glog.Debug("main", TOTAL_TOPICS, len(hub.subscriptions))
		}

		close(client.send)
	}
}
