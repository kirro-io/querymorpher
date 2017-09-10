# querymorpher
[![Build Status](https://travis-ci.org/Kirro-io/querymorpher.svg?branch=master)](https://travis-ci.org/Kirro-io/querymorpher)
[![codecov](https://codecov.io/gh/Kirro-io/querymorpher/branch/master/graph/badge.svg)](https://codecov.io/gh/Kirro-io/querymorpher)
[![Go Report Card](https://goreportcard.com/badge/github.com/kirro-io/querymorpher)](https://goreportcard.com/report/github.com/kirro-io/querymorpher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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
		q, err := querymorpher.QueryFromRequest(r)

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

## Ordering results

For ascending ordering use `order_by=<field_name>`:

```
?age__gt=18&order_by=age
```

For descending ordering insert minus sign before field name `order_by=-<field_name>`:

```
?age__gt=18&order_by=-age
```

## Limiting results

Simply use `limit=<count>` query param.

```
?age__gt=18limit=1
```
