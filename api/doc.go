// Package api is a client for the Uspacy REST API.
//
// Create one client per API host with [New] and reuse it: it is safe for concurrent use,
// and concurrent calls that hit an expired token share a single token refresh.
//
//	client := api.New(accessToken, refreshToken, "https://example.uspacy.ua",
//		api.WithMaxRetries(5),
//	)
//
// The types of requests and responses live in the sibling packages (crm, task, user, …).
//
// # Context
//
// Every method that sends a request takes a [context.Context] first. Cancelling it, or
// letting its deadline pass, stops the request in flight, any wait before a retry and a
// token refresh the call is waiting on.
//
// # Request options
//
// Every such method also takes [RequestOption] values last, for example [WithHeader]:
//
//	leads, err := client.GetLeads(ctx, params, api.WithHeader("X-Request-Id", reqID))
//
// # Errors
//
// A final non-2xx response is returned as an [*HTTPError] carrying the status code and the
// response body; use [errors.As] to inspect it. A failed token refresh after a 401 is a
// 401 [*HTTPError] whose Err is the refresh error. When ctx ends, the error wraps ctx.Err(),
// so [errors.Is] with [context.Canceled] or [context.DeadlineExceeded] works.
//
// # Retries
//
// A 429 is retried after its Retry-After delay. A 5xx response or a transport error is
// retried with exponential backoff for idempotent methods (GET, HEAD, OPTIONS, PUT, DELETE).
// POST and PATCH are retried only when the connection was never established, so a write
// the server may already have applied is not sent twice. A 401 triggers one token refresh
// and a retry. Tune with [WithMaxRetries] and [WithRetryBackoff].
//
// # Tokens
//
// The client authenticates with the access token and refreshes it with the refresh
// token on a 401. [Uspacy.Tokens] returns the current pair, for callers that store it.
//
// # Upgrading from v1
//
// See MIGRATION.md in the repository root:
// https://github.com/Uspacy/uspacy-go-sdk/blob/main/MIGRATION.md
package api
