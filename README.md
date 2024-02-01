WebSocket Communication System
=================================

Introduction
------------

This project is a Go-based WebSocket communication system designed to handle real-time data exchange between producers and consumers.

The system allows multiple producers to send data to a centralized hub, which then broadcasts this data to subscribed consumers. This setup is ideal for applications requiring real-time data feeds, such as chat applications, live data updates, and more.

__NOTE__: This is only for illustration and does not include any authentication. This was also written in 2018 and is probably very outdated.

Getting Started
---------------

### Prerequisites

*   Go (version 1.x or later)
*   [WebSocket](https://github.com/gorilla/websocket) Go package

### Installation

1.  Clone the repository to your local machine.

    bashCopy code

    `git clone [URL to your repository]`

2.  Navigate to the project directory.

    bashCopy code

    `cd [your directory]`

3.  Install necessary Go packages.

    bashCopy code

    `go get`


### Running the Server

1.  Start the server by running:

    bashCopy code

    `go run .`

    The server listens on port 8888 by default but can be configured using the `-addr` flag.

### Configuration

*   The server's log level can be set to `info`, `debug`, or `trace` by passing it as the first command-line argument.
*   You can change the default port by using the `-addr` flag followed by `:<port_number>`.

Testing
-------

### Simulating Producers

Producers can send data to the server using WebSocket connections. Here's a basic JavaScript example to simulate a producer in JavaScript:

```javascript
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8888/pub/producer1');

ws.on('open', function open() {
    ws.send(JSON.stringify({Topic: "testTopic", Data: "Hello World" }));
});

ws.on('message', function incoming(data) {
    console.log(data);
});
```

### Simulating Consumers

Consumers subscribe to topics to receive data. Here's an example script for a consumer:

```javascript
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8888/sub/testTopic');

ws.on('message', function incoming(data) {
    console.log("Received:", data);
});
```