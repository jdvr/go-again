# Go Again
---

A simple and configurable retry library for go, with exponential backoff, and constant delay support out of the box.
Inspired by [backoff](https://github.com/cenkalti/backoff).

## Features

- Configurable delay calculation algorithm
- Support for exponential backoff and constant delay out of the box
- Support for generics
- Simple and clean interface

### There are two main concepts:
- Retry: Given an operation and a ticks calculator keeps retrying until either permanent error or timeout happen
- TicksCalculator: Provide delay for retryer to wait between retries


## Examples:

### Call an API using exponential backoff

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jdvr/go-again"
)


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiResponse, err := again.Retry[*http.Response](ctx, func(ctx context.Context) (*http.Response, error) {
		fmt.Println("Running Operation")

		resp, err := http.DefaultClient.Get("https://sameflaky.api/path")
		if err != nil {
			// operation will be retried
			return nil, err
		}

		if resp.StatusCode == http.StatusForbidden {
			// no more retries
			return nil, again.NewPermanentError(errors.New("no retry, permanent error"))
		}

		if resp.StatusCode > 400 {
			return nil, errors.New("this will be retry")
		}

		// do whatever you need with a valid response ...

		return resp, nil // no retry
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished with response %v\n", apiResponse)
}

```
## Call database keeping a constant delay
```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jdvr/go-again"
)

type Result struct {
	value string
}

type dbOperation struct {
	repo Repo
}

func(db dbOperation) Run(ctx context.Context) (Result, error) {
	fmt.Println("Running Operation")

	resp, err := db.repo.GetAll()
	if err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			return Result{}, again.NewPermanentError(fmt.Errorf("no retry, permanent error: %w", err))

		}
		// operation will be retried
		return Result{}, err
	}

	// do whatever you need with a valid response ...

	return resp, nil // no retry
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	constantRetryer := again.WithConstantDelay[Result](15 * time.Millisecond, 30 * time.Second)

	result, err := constantRetryer.Retry(ctx, dbOperation{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished %v\n", result)
}
```


## Test

`make test`

## Next steps:

- Get rid of manual mocks
- Improve flaky testing (depends on system clock)
- Configure the linter