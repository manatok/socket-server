package main

const (
	APP_STARTED            = "AppStarted: (%#v)"
	REGISTERING_CLIENT     = "RegisteringClient: (%#v)"
	UNREGISTERING_CLIENT   = "UNRegisteringClient: (%#v)"
	CONNECTING_PRODUCER    = "ConnectingProducer: (%#v)"
	CONNECTING_CLIENT      = "ConnectingClient: (%#v)"
	SOCKET_CREATE_ERROR    = "SocketCreateError: (%s)"
	WEBSERVER_CREATE_ERROR = "WebserverCreateError: (%s)"
	STARTING_HUB           = "StartingHub: ()"
	UNREGISTERING_TOPIC    = "UnregisteringTopic: (%s)"
	TOTAL_TOPICS           = "TotalTopics: (%d)"
	TOTAL_TOPIC_CLIENTS    = "TotalClientTopics: (%s %d)"
	MESSAGE_BROADCAST      = "MessageBroadcast: (%s %d)"
	CLIENT_BUFFER_FULL     = "ClientBufferFull: (%#v)"
	SENDING_MESSAGE        = "SendingMessage: (%s)"
)
