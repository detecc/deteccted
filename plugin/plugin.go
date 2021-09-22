package plugin

type Handler interface {
	// Execute method is called when a client receives a command from the server.
	// The arguments of the method is Payload's data attribute, incoming from the server. The response should be data, ready to
	// be sent back to the server for processing and an error, if one occurred.
	Execute(args interface{}) (interface{}, error)
}