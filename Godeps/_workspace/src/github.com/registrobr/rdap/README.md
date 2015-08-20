RDAP
====

RDAP (Registration Data Access Protocol) is a library to be used in clients and
servers to make the life easier when building requests and responses. You will
find all RDAP protocol types in the protocol package and can use the clients to
build your own client tool.

Implements the RFCs:
  * 7480 - HTTP Usage in the Registration Data Access Protocol (RDAP)
  * 7482 - Registration Data Access Protocol (RDAP) Query Format
  * 7483 - JSON Responses for the Registration Data Access Protocol (RDAP)
  * 7484 - Finding the Authoritative Registration Data (RDAP) Service

Also support the extensions:
  * NIC.br RDAP extension

Usage
-----

Download the project with:

```
go get github.com/registrobr/rdap
```

And build a program like bellow for direct RDAP server requests:

```go
package main

import (
	"fmt"

	"github.com/registrobr/rdap"
)

func main() {
	c := rdap.NewClient([]string{"https://rdap.beta.registro.br"})

	d, err := c.Query("nic.br", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v", d)
}
```

You can also try with bootstrap support:

```go
package main

import (
	"fmt"

	"github.com/registrobr/rdap"
)

func main() {
	c := rdap.NewClient(nil)

	ipnetwork, err := c.Query("214.1.2.3", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v", ipnetwork)
}
```

For advanced users you probably want to reuse the HTTP client and add a cache
layer:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/registrobr/rdap"
)

func main() {
	var httpClient http.Client

	cacheDetector := rdap.CacheDetector(func(resp *http.Response) bool {
		return resp.Header.Get("X-From-Cache") == "1"
	})

	c := rdap.Client{
		Transport: rdap.NewBootstrapFetcher(&httpClient, rdap.IANABootstrap, cacheDetector),
	}

	ipnetwork, err := c.Query("214.1.2.3", http.Header{
		"X-Forwarded-For": []string{"127.0.0.1"},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v", ipnetwork)
}
```

An example of usage can be found in the project:
[https://github.com/registrobr/rdap-client](https://github.com/registrobr/rdap-client)
