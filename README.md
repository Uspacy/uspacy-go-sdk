# uspacy-go-sdk

Go client for the Uspacy REST API: CRM, tasks, users, groups,
files, activities, notifications and more.

[![Go Reference](https://pkg.go.dev/badge/github.com/Uspacy/uspacy-go-sdk/v2.svg)](https://pkg.go.dev/github.com/Uspacy/uspacy-go-sdk/v2/api)

## Install

```sh
go get github.com/Uspacy/uspacy-go-sdk/v2
```

Requires Go 1.23 or newer.

## Usage

```go
import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/api"
)

client := api.New(accessToken, refreshToken, "https://example.uspacy.ua")

leads, err := client.GetLeads(ctx, url.Values{"page": {"1"}},
	api.WithHeader("X-Request-Id", reqID),
)

var httpErr *api.HTTPError
if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusForbidden {
	// httpErr.Body holds the response body
}
```

- Every call takes a `context.Context` first and `api.RequestOption` values last.
- Create one client per API host and reuse it: it is safe for concurrent use and refreshes
  the access token on a 401.
- Non-2xx responses are returned as `*api.HTTPError`. Rate limits (429), server errors and
  network failures are retried; POST and PATCH are not re-sent once they may have reached
  the server.

The [package documentation](https://pkg.go.dev/github.com/Uspacy/uspacy-go-sdk/v2/api)
covers retries, errors and tokens in detail, with examples.

## Upgrading from v1

v2 changes every method signature. See [MIGRATION.md](MIGRATION.md).

## Development

```sh
go vet ./...
go test -race ./...
```

CI runs these, plus gofmt, `go mod tidy` and staticcheck, on every pull request against
Go 1.23 and the latest stable release.
