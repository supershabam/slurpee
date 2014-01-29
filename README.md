# slurpee

I subscribe to a redis server and give you a channel of bytes. I handle disconnects and reconnects for you. Just think about all those tasty bytes you could be processing.

## api

```go
import(
	"github.com/supershabam/slurpee"
)

var s, err := slurpee.NewSlurpee("redis://u:password@localhost", "my-channel")
if err != nil {
	panic(err)
}
for b := range s.Channel() {
	fmt.Println("got bytes! %v", b)	  
	// handle(b)
}
```

This example merely prints the bytes out to the console, but try running it and restarting your redis server.

## license

MIT

