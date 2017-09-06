# querymorpher

[![Go Report Card](https://goreportcard.com/badge/github.com/kirro-io/querymorpher)](https://goreportcard.com/report/github.com/kirro-io/querymorpher)

Simple package that converts `url.Values` to sql like query.

## Installation

```
$ go get github.com/kirro-io/querymorpher
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kirro-io/querymorpher"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q, err := querymorpher.Transform(r.URL.Query())

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Fprintf(w, q)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

```shell
$ curl "http://localhost:8080?age>=42&name=John&surname='Doe'"
age >= 42 AND name = 'John' AND surname = 'Doe'
```

## Supported operators

suffix | operator
-------|---------
__gt | >
__gte | >=
__lt | <
__lte | <=
__neq | !=
