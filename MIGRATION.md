# Migrating from v1 to v2

v2 makes every network call cancellable with `context.Context`, replaces header maps with
request options, reports HTTP failures as a typed error and makes retries safe for writes.
Most changes are caught by the compiler; the ones it cannot catch are marked **Runtime**
below. Read those even if your code compiles.

## Checklist

1. Switch to the `/v2` module path and Go 1.23+ ([details](#module-path-and-go-version)).
2. Pass a `context.Context` as the first argument of every SDK call ([details](#contextcontext-on-every-call)).
3. Replace `map[string]string` header arguments with `api.WithHeader` / `api.WithHeaders` ([details](#request-headers)).
4. Replace `client.RefreshToken` with `client.Tokens()` and `TokenRefresh()` with `TokenRefresh(ctx)` ([details](#tokens)).
5. Rename `MassDeletionActivies` and `GetFunnelStageDyId` ([details](#renamed-methods)).
   Pass `nil` or a `url.Values` to `GetGroups` and `GetTaskById` ([details](#optional-query-parameters)).
6. **Runtime:** check code that matches on error text ([details](#errors)).
7. **Runtime:** check code that relies on POST/PATCH being retried ([details](#retries)).
8. **Runtime:** check calls to `DeleteFilesByEntityId` ([details](#other-behavior-changes)).

## Module path and Go version

The module is now `github.com/Uspacy/uspacy-go-sdk/v2` and requires Go 1.23 or newer.

```sh
go get github.com/Uspacy/uspacy-go-sdk/v2@latest
# rewrite imports; the lookahead keeps an already migrated path from becoming /v2/v2
find . -name '*.go' -not -path './vendor/*' -exec \
  perl -pi -e 's#github\.com/Uspacy/uspacy-go-sdk/(?!v2/)#github.com/Uspacy/uspacy-go-sdk/v2/#g' {} +
go mod tidy
go mod vendor   # only if the service vendors its dependencies
```

## `context.Context` on every call

Every method that performs a request takes `ctx` as its first argument. Cancelling `ctx`, or
letting its deadline pass, stops the request in flight, the wait before a retry and a token
refresh the call is waiting on.

```go
// v1
leads, err := client.GetLeads(params)

// v2
leads, err := client.GetLeads(ctx, params)
```

Pass the incoming request's context in HTTP handlers, the job's context in workers, and
`context.Background()` where nothing better exists. A ctx error can be checked with
`errors.Is(err, context.DeadlineExceeded)` or `errors.Is(err, context.Canceled)`.

## Request headers

Every method that sends a request now takes `opts ...api.RequestOption` as its last
parameter. Methods that accepted `headers ...map[string]string` take options instead.

```go
// v1
id, code, err := client.CreateEntity("leads", data, map[string]string{"X-Request-Id": reqID})

// v2
id, code, err := client.CreateEntity(ctx, "leads", data, api.WithHeader("X-Request-Id", reqID))
// or, with an existing map:
id, code, err := client.CreateEntity(ctx, "leads", data, api.WithHeaders(headers))
```

- Empty values are still ignored, as in v1.
- A header replaces the default of the same name, including `Authorization`. In v1 it was
  added next to it, so the request carried two values.
- In v1, passing any header dropped `Content-Type: application/json` from JSON requests.
  v2 keeps it.
- The `Content-Type` of file uploads and form requests cannot be overridden: the body is
  encoded for that exact type.
- If your code stores SDK methods in variables or mocks the client behind an interface, add
  `opts ...api.RequestOption` to those signatures as well.

## Tokens

```go
// v1
newToken, err := client.TokenRefresh()
refresh := client.RefreshToken

// v2
newToken, err := client.TokenRefresh(ctx)
_, refresh := client.Tokens()
```

- The `RefreshToken` field is gone. `Tokens()` returns the access and refresh token as one
  consistent pair; reading them separately could mix tokens from two different refreshes.
- **Runtime:** the refresh request now authenticates with the **refresh** token
  (`api.New(token, refresh, host)`). v1 sent the access token, so once it had expired the
  automatic refresh on a 401 could not succeed. `api.New(token, token, host)` and
  `api.New(token, "", host)` behave exactly as in v1.
- Concurrent calls on the same client that get a 401 share a single refresh request. A call
  whose token was already refreshed by another one retries with the new token instead of
  refreshing again. This works per client: a service that creates a new client for every
  call gets neither, so prefer one long-lived client per host.
- If the service stores tokens (for example in a database), read `Tokens()` after the
  calls: the client may have refreshed them on its own after a 401.

## Renamed methods

| v1 | v2 |
|---|---|
| `MassDeletionActivies` | `MassDeleteActivities` |
| `GetFunnelStageDyId` | `GetFunnelStageById` |

## Optional query parameters

`GetGroups` and `GetTaskById` took query parameters as `params ...url.Values`. Go allows only
one variadic parameter, and it is now the request options, so `params` is a plain
`url.Values`, as in `GetLeads` and the other list methods. Pass `nil` when there are none.

```go
// v1
groups, err := client.GetGroups()
groups, err := client.GetGroups(params)
task, err := client.GetTaskById(id)

// v2
groups, err := client.GetGroups(ctx, nil)
groups, err := client.GetGroups(ctx, params)
task, err := client.GetTaskById(ctx, id, nil)
```

`GetGroups` used to merge several `url.Values` passed to it; merge them before the call.

## Errors

**Runtime.** Code that only checks `err != nil` needs no change. Code that inspects errors
should check the cases below.

### Non-2xx responses are `*api.HTTPError`

```go
var httpErr *api.HTTPError
if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
	// httpErr.Body holds the response body
}
```

The text starts exactly as in v1, so `strings.Contains(err.Error(), "status code: 422")`
keeps working: `request failed: [POST] <url>, status code: 422, response: <body>`.

| Case | v1 | v2 |
|---|---|---|
| Final 3xx response | success, `err == nil` | `*HTTPError` |
| 401 and the token refresh fails | the refresh error as is | `*HTTPError` with `StatusCode` 401 and `Err` set to the refresh error; the text ends with `, cause: <refresh error>`, so text checks such as `status code: 403` + `unauthenticated` still match |
| `ctx` ends after a 429/5xx | (no ctx) | that `*HTTPError`, with `Err` set to the ctx error |
| `ctx` ends otherwise | (no ctx) | `request aborted: <ctx error>` |
| No HTTP response on any attempt | `request failed after 3 retries:` … | `request failed (attempts: N):` … |

The last row changes text: search for `after 3 retries` or `retries:` in code that parses
these errors. In v2, N is the number of requests actually sent; it is 1 for a POST/PATCH
that timed out (see [Retries](#retries)).

## Retries

**Runtime.** v2 no longer re-sends a write the server may already have applied.

| | v1 | v2 |
|---|---|---|
| Attempts | 3 | 3, configurable with `api.WithMaxRetries(n)` |
| Delay between attempts | about 3 s, then about 5 s | doubles from 3 s up to 30 s, randomised by up to half; configurable with `api.WithRetryBackoff(base, max)` |
| 429 | retried after `Retry-After` | retried after `Retry-After` (capped at 5 min); without the header, after the normal delay |
| 5xx | retried for every method | retried for GET, HEAD, OPTIONS, PUT, DELETE only |
| Timeout or connection reset | retried for every method | retried for GET, HEAD, OPTIONS, PUT, DELETE; for POST/PATCH only on a dial or DNS failure, where the connection was never established (other failures before sending, such as a TLS handshake error, are not retried) |
| 401 | one token refresh, then the request again | the same |

If the service wraps SDK calls in its own retry loop, that loop still re-sends POST/PATCH
requests after a timeout and can create duplicates. Consider retrying writes there only when
the error shows the request never reached the server.

To use your own HTTP client (proxy, tracing, other timeout), pass `api.WithHTTPClient(c)`.
The default client has a 30 s timeout per attempt.

## Other behavior changes

- **Runtime:** `DeleteFilesByEntityId` now sends `?entityId=<id>&entityType=<type>`. v1 sent
  `files?<type>&<id>` without parameter names. Check what the call does against your portal
  before relying on it: it now passes the filter the endpoint expects.
- Methods that return a list unwrapped from the response (`GetFields`, `GetTaskFields`,
  `GetAllFunnelStages`, …) return `nil` together with a JSON decode error, not a partially
  decoded list.
- `GetTasksWithFilters` returns an error if the response's `data` is not a list of objects.
  v1 silently returned no tasks, or `nil` entries for items that were not objects.

## New in v2

- `GetEntity(ctx, entityType, id)` returns one CRM record as raw JSON.
- `crm.Field.Dependency` describes a dependent list field; `crm.Value.Active` tells an
  explicit `false` from an absent key (`nil`).
- Client options: `api.WithMaxRetries`, `api.WithRetryBackoff`, `api.WithHTTPClient`.
