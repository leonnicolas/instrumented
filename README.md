# instrumented

## client

Instrument a net/http client as easy as:

```go
package main

import (
	ic "github.com/leonnicolas/instrumented/client"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	r := prometheus.NewRegistry()
	c := ic.NewInstrumentedClient(prometheus.WrapRegistererWithPrefix("prefix_", r))
	...
}
```

or instrument an existing client with:

```go
hc := &http.Client{}

c := ic.InstrumentClient(hc, prometheus.WrapRegistererWithPrefix("prefix_", r))
```
