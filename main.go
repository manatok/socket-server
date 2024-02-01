package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8888", "http service address")

func init() {
	logLevel := glog.INFO

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "info":
			logLevel = glog.INFO
		case "debug":
			logLevel = glog.DEBUG
		case "trace":
			logLevel = glog.TRACE
		}
	}

	glog.ConfigureLoggers(logLevel)
}

func main() {
	glog.Info("main", APP_STARTED, os.Args)

	flag.Parse()

	// The hub holds all of the subscribers
	hub := newHub()
	go hub.run()

	// Allow cross-domain requests.
	corsObj := handlers.AllowedOrigins([]string{"*"})

	rtr := mux.NewRouter()

	// The endpoint that producers would publish on.
	rtr.HandleFunc("/pub/{token:[a-zA-Z0-9_]+}", func(w http.ResponseWriter, r *http.Request) {
		connectProducer(hub, w, r)
	}).Methods("GET")

	rtr.HandleFunc("/sub/{topic:[a-zA-Z0-9_]+}", func(w http.ResponseWriter, r *http.Request) {
		connectClient(hub, w, r)
	}).Methods("GET")

	// Start the webserver
	err := http.ListenAndServe(*addr, handlers.CORS(corsObj)(rtr))

	if err != nil {
		glog.Critical("main", WEBSERVER_CREATE_ERROR, err)
	}
}

func connectProducer(hub *Hub, w http.ResponseWriter, r *http.Request) {
	glog.Debug("main", CONNECTING_PRODUCER, mux.Vars(r))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error("main", SOCKET_CREATE_ERROR, err)
		return
	}

	producer := &Producer{hub: hub, conn: conn}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go producer.readPump()
}

func connectClient(hub *Hub, w http.ResponseWriter, r *http.Request) {
	glog.Debug("main", CONNECTING_CLIENT, mux.Vars(r))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error("main", SOCKET_CREATE_ERROR, err)
		return
	}

	params := mux.Vars(r)
	topic := params["topic"]

	// This buffer size will need to be tweaked
	client := &Client{topic: topic, conn: conn, send: make(chan []byte, 20)}
	hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}
