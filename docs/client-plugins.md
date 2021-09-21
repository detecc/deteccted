# Client plugins

Client plugins must implement the `Handler` interface defined in the `plugin` package:

```golang
package main
import github.com/detecc/detecctor/shared
type Handler interface {
// Execute method is called when a client receives a command from the server. As an argument, it gets the data attribute from the Payload. 
//The execute method returns the Payload ready to be sent as a reply to the server's command request. 
Execute(args interface{}) (shared.Payload, error)
}
```

