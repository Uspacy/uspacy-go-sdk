package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

func TestBuildURL(t *testing.T) {
	us := &Uspacy{mainHost: "https://example.uspacy.ua"}

	cases := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			"joins parts with single slashes",
			[]string{"crm/v1/", "entities/contacts/", "5"},
			"https://example.uspacy.ua/crm/v1/entities/contacts/5",
		},
		{
			"trims duplicate slashes",
			[]string{"/crm/v1/", "/entities/contacts/"},
			"https://example.uspacy.ua/crm/v1/entities/contacts",
		},
		{
			"skips empty parts",
			[]string{"crm/v1", "", "entities/deals"},
			"https://example.uspacy.ua/crm/v1/entities/deals",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := us.buildURL(tc.parts...); got != tc.want {
				t.Errorf("buildURL(%v) = %q, want %q", tc.parts, got, tc.want)
			}
		})
	}
}

// Regression for PatchEntity/EntityMassEdit building URLs without a slash
// between the entity and the trailing segment: buildURL trims the trailing
// slash that crm.EntityUrl carries, so appending the id with plain string
// concatenation produced /crm/v1/entities/contacts5 (matched by the backend
// as the collection route and rejected with 405). The trailing segment must
// be passed to buildURL as its own part.
func TestPatchEntityURLHasSlashBeforeID(t *testing.T) {
	us := &Uspacy{mainHost: "https://example.uspacy.ua"}

	got := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, "contacts"), "5")
	want := "https://example.uspacy.ua/crm/v1/entities/contacts/5"
	if got != want {
		t.Errorf("PatchEntity URL = %q, want %q", got, want)
	}

	got = us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, "contacts"), "mass_edit")
	want = "https://example.uspacy.ua/crm/v1/entities/contacts/mass_edit"
	if got != want {
		t.Errorf("EntityMassEdit URL = %q, want %q", got, want)
	}
}

func TestGetEntityProductList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/crm/v1/static/entity-product-lists" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/crm/v1/static/entity-product-lists")
		}
		if got := r.URL.Query().Get("entity_type"); got != "deals" {
			t.Errorf("entity_type = %q, want %q", got, "deals")
		}
		if got := r.URL.Query().Get("entity_id"); got != "1072" {
			t.Errorf("entity_id = %q, want %q", got, "1072")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":5,"entity_type":"deals","entity_id":1072,"is_automatic_calculation":0,"amount_before_discount_and_tax":977500,"amount_discount":-534653,"amount_tax":284609,"amount_before_tax":1423044,"amount_total":1707653,"list_products":[{"id":8,"title":"Вентилятор 12","price":1173000,"currency":"UAH","quantity":1,"price_type_id":null,"measurement_unit_abbr":"pcs","discount_value":-4558,"discount_type":"relative","discount_price":0,"tax_rate":20,"is_tax_included":1,"amount":531956,"created_at":1789130148,"updated_at":1789130148,"product":null}]}`)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	productList, err := us.GetEntityProductList("deals", 1072)
	if err != nil {
		t.Fatalf("GetEntityProductList() error = %v", err)
	}
	if productList.ID != 5 || productList.EntityType != "deals" || productList.EntityID != 1072 {
		t.Errorf("GetEntityProductList() = %+v", productList)
	}
	if len(productList.ListProducts) != 1 || productList.ListProducts[0].Title != "Вентилятор 12" {
		t.Errorf("ListProducts = %+v", productList.ListProducts)
	}
}

func TestCreateEntityListProduct(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/crm/v1/static/list-products" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/crm/v1/static/list-products")
		}
		var request crm.CreateEntityListProductRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Title != "кк" || request.EntityProductListID != 5 || request.TaxRate != nil {
			t.Errorf("request body = %+v", request)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":9,"title":"кк","price":0,"currency":"UAH","quantity":1,"price_type_id":null,"measurement_unit_abbr":"pcs","discount_value":0,"discount_type":"relative","discount_price":0,"tax_rate":null,"is_tax_included":0,"amount":0,"created_at":1789142588,"updated_at":1789142588,"product":null}`)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	product, err := us.CreateEntityListProduct(crm.CreateEntityListProductRequest{
		Currency:            "UAH",
		DiscountType:        "relative",
		MeasurementUnitAbbr: "pcs",
		Quantity:            1,
		Title:               "кк",
		EntityProductListID: 5,
	})
	if err != nil {
		t.Fatalf("CreateEntityListProduct() error = %v", err)
	}
	if product.ID != 9 || product.Title != "кк" || product.TaxRate != nil {
		t.Errorf("CreateEntityListProduct() = %+v", product)
	}
}

func TestDeleteEntityListProduct(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodDelete)
		}
		if r.URL.Path != "/crm/v1/static/list-products/9" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/crm/v1/static/list-products/9")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	statusCode, err := us.DeleteEntityListProduct(9)
	if err != nil {
		t.Fatalf("DeleteEntityListProduct() error = %v", err)
	}
	if statusCode != http.StatusNoContent {
		t.Errorf("DeleteEntityListProduct() status = %d, want %d", statusCode, http.StatusNoContent)
	}
}

// --- Characterization tests for doRawInternal/doRawInternalCtx ---
//
// These pin today's non-context retry and error behaviour before the context-aware
// refactor, so a regression in doRawInternalCtx (which every existing method now
// runs through via context.Background()) shows up here. The 5xx case sleeps through
// two real backoffs (about 8s); that cost is accepted for a characterization test.

func TestDoRawRetries5xxButNot4xx(t *testing.T) {
	cases := []struct {
		name         string
		status       int
		wantRequests int32
	}{
		{"5xx is retried up to defaultRetries", http.StatusInternalServerError, defaultRetries},
		{"4xx is not retried", http.StatusBadRequest, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requests int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				w.WriteHeader(tc.status)
			}))
			defer server.Close()

			us := New("token", "", server.URL)
			_, statusCode, err := us.doRaw(us.buildURL("resource"), http.MethodGet, nil, nil)

			if err == nil {
				t.Error("doRaw() error = nil, want a non-2xx error")
			}
			if statusCode != tc.status {
				t.Errorf("statusCode = %d, want %d", statusCode, tc.status)
			}
			if got := atomic.LoadInt32(&requests); got != tc.wantRequests {
				t.Errorf("server received %d requests, want %d", got, tc.wantRequests)
			}
		})
	}
}

func TestDoRaw401RefreshFailureReturnsRawError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// "token" has no dots, so UnmarshalTokenData rejects it before TokenRefresh ever
	// makes a network call. That failure must reach the caller unwrapped.
	us := New("token", "", server.URL)
	_, statusCode, err := us.doRaw(us.buildURL("resource"), http.MethodGet, nil, nil)

	if statusCode != http.StatusUnauthorized {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusUnauthorized)
	}
	wantErr := "invalid JWT token format: expected 3 parts, got 1"
	if err == nil || err.Error() != wantErr {
		t.Errorf("err = %v, want %q", err, wantErr)
	}
	if httpErr, ok := err.(*HTTPError); ok {
		t.Errorf("err = %#v, want the raw refresh error, not *HTTPError", httpErr)
	}
}

// The old path must return the refresh failure's exact value and type, not a wrapper around
// it: doRawInternal strips the *refreshFailedError marker doRawInternalCtx produces, so
// errors.Unwrap must find nothing further, and a type assertion (not just errors.As) must
// reach the refresh's own error directly.
func TestDoRawOldPathReturnsRefreshErrorUnwrapped(t *testing.T) {
	t.Run("malformed token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		us := New("token", "", server.URL)
		_, _, err := us.doRaw(us.buildURL("resource"), http.MethodGet, nil, nil)
		if err == nil || errors.Unwrap(err) != nil {
			t.Errorf("err = %v (Unwrap = %v), want the raw refresh error with no wrapper", err, errors.Unwrap(err))
		}
	})
	t.Run("refresh answers 403", func(t *testing.T) {
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/auth/refresh_token") {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()
		us := New(testJWT(t, strings.TrimPrefix(srv.URL, "https://")), "", srv.URL)
		us.client = srv.Client()
		_, err := us.GetFields("leads")
		he, ok := err.(*HTTPError) // a type assertion, not errors.As: the value itself
		if !ok || he.StatusCode != http.StatusForbidden || he.Method != http.MethodPost {
			t.Errorf("GetFields() err = %T %v, want the refresh's own *HTTPError (POST 403) unwrapped", err, err)
		}
	})
}

func TestDoRawFinalNon2xxErrorText(t *testing.T) {
	const respBody = `{"message":"not found"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, respBody)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	url := us.buildURL("resource")
	body, statusCode, err := us.doRaw(url, http.MethodGet, nil, nil)

	if statusCode != http.StatusNotFound {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusNotFound)
	}
	if string(body) != respBody {
		t.Errorf("body = %q, want %q", body, respBody)
	}
	wantErr := fmt.Sprintf("request failed: [%s] %s, status code: %d, response: %s", http.MethodGet, url, http.StatusNotFound, respBody)
	if err == nil || err.Error() != wantErr {
		t.Errorf("err = %v, want %q", err, wantErr)
	}
}

func TestHTTPErrorUnwrap(t *testing.T) {
	cause := errors.New("refresh failed")
	err := &HTTPError{Method: http.MethodGet, URL: "https://example.uspacy.ua/x", StatusCode: http.StatusUnauthorized, Err: cause}

	if !errors.Is(err, cause) {
		t.Error("errors.Is(err, cause) = false, want true")
	}
	wantText := "request failed: [GET] https://example.uspacy.ua/x, status code: 401, response: "
	if err.Error() != wantText {
		t.Errorf("Error() = %q, want %q", err.Error(), wantText)
	}
}

func TestFinalizeCtxError(t *testing.T) {
	const method, url = http.MethodGet, "https://example.uspacy.ua/x"

	if got := finalizeCtxError(method, url, http.StatusOK, nil, nil); got != nil {
		t.Errorf("finalizeCtxError(200, nil, nil) = %v, want nil", got)
	}

	already := &HTTPError{Method: method, URL: url, StatusCode: http.StatusInternalServerError}
	if got := finalizeCtxError(method, url, http.StatusInternalServerError, nil, already); got != already {
		t.Error("finalizeCtxError should return an existing *HTTPError unchanged, not re-wrap it")
	}

	// A context error (or any other non-HTTPError, non-refreshFailedError error) must pass
	// through unchanged, whatever statusCode says: dropping asHTTPError's statusCode >= 400
	// wrapping rule means finalizeCtxError never wraps an arbitrary error as *HTTPError.
	ctxErr := fmt.Errorf("request aborted: %w", context.DeadlineExceeded)
	if got := finalizeCtxError(method, url, 0, nil, ctxErr); got != ctxErr {
		t.Errorf("finalizeCtxError(0, ctxErr) = %v, want ctxErr unchanged", got)
	}

	// A refresh failure marker becomes a 401 *HTTPError whose Err is the refresh error,
	// regardless of what statusCode doRawInternalCtx returned alongside it.
	refreshErr := errors.New("refresh boom")
	got := finalizeCtxError(method, url, http.StatusUnauthorized, nil, &refreshFailedError{refreshErr})
	var httpErr *HTTPError
	if !errors.As(got, &httpErr) {
		t.Fatalf("finalizeCtxError(401, refreshFailedError) = %v, want *HTTPError", got)
	}
	if httpErr.Method != method || httpErr.URL != url || httpErr.StatusCode != http.StatusUnauthorized || httpErr.Err != refreshErr {
		t.Errorf("finalizeCtxError(401, refreshFailedError) = %+v, want Method=%q URL=%q StatusCode=401 Err=refreshErr", httpErr, method, url)
	}

	// A final 3xx is success everywhere else (doRawInternalCtx, every non-context method),
	// but finalizeCtxError reports it as an *HTTPError with its body for GetFieldsCtx and
	// GetEntityCtx.
	got = finalizeCtxError(method, url, http.StatusFound, []byte("moved"), nil)
	if !errors.As(got, &httpErr) {
		t.Fatalf("finalizeCtxError(302, nil) = %v, want *HTTPError", got)
	}
	if httpErr.StatusCode != http.StatusFound || string(httpErr.Body) != "moved" {
		t.Errorf("finalizeCtxError(302, nil) = %+v, want StatusCode=302 Body=%q", httpErr, "moved")
	}
}

func TestRefreshFailedError(t *testing.T) {
	cause := errors.New("refresh boom")
	err := &refreshFailedError{cause}

	if err.Error() != cause.Error() {
		t.Errorf("Error() = %q, want %q", err.Error(), cause.Error())
	}
	if !errors.Is(err, cause) {
		t.Error("errors.Is(err, cause) = false, want true")
	}
}

// requestFailedAfterRetries builds the error for doRawInternalCtx's post-loop path, taken
// when every attempt fails at the transport level (statusCode stays 0). Exercising that
// exact path through doRawInternalCtx itself would need two real backoff sleeps
// (defaultRetries=3, with jittered ~3s/~5s delays) to run to completion without ctx firing,
// and then ctx to end in the narrow window right after the third attempt's immediate
// failure — the only attempt with no following sleep to catch a ctx that ends during it.
// That window can't be hit deterministically (the two backoffs jitter independently across
// a ~2s range each), so this test drives the extracted helper directly instead.
func TestRequestFailedAfterRetries(t *testing.T) {
	logs := map[string]*errorLog{
		"dial tcp: connection refused": {message: "dial tcp: connection refused", attempts: []int{1, 2, 3}},
	}
	wantDetails := "Error 'dial tcp: connection refused' on attempts: [1 2 3]\n"

	t.Run("ctx.Err() nil keeps the pre-context text and type", func(t *testing.T) {
		err := requestFailedAfterRetries(context.Background(), defaultRetries, logs)
		wantText := fmt.Sprintf("request failed after %d retries:\n%s", defaultRetries, wantDetails)
		if err.Error() != wantText {
			t.Errorf("Error() = %q, want %q", err.Error(), wantText)
		}
		if errors.Unwrap(err) != nil {
			t.Errorf("Unwrap() = %v, want nil: no %%w chain when ctx.Err() is nil", errors.Unwrap(err))
		}
	})

	t.Run("canceled ctx is wrapped with %w", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := requestFailedAfterRetries(ctx, defaultRetries, logs)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled in its chain", err)
		}
		wantPrefix := fmt.Sprintf("request failed after %d retries:\n%scontext error: ", defaultRetries, wantDetails)
		if !strings.HasPrefix(err.Error(), wantPrefix) {
			t.Errorf("Error() = %q, want prefix %q", err.Error(), wantPrefix)
		}
	})

	t.Run("expired deadline is reachable via errors.Is", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()
		<-ctx.Done()

		err := requestFailedAfterRetries(ctx, defaultRetries, logs)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
	})
}

// --- Context-path tests for doRawInternalCtx ---

// testJWT builds a syntactically valid but unsigned JWT whose payload sets the
// "domain" claim, which is all UnmarshalTokenData/tokenRefreshCtx need from it.
func testJWT(t *testing.T, domain string) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"domain":%q}`, domain)))
	return header + "." + payload + ".sig"
}

func TestDoRawInternalCtxRetryWaitCutShortBy502(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "bad gateway")
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, statusCode, err := us.doRawInternalCtx(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, true)
	elapsed := time.Since(start)

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err = %v, want *HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusBadGateway {
		t.Errorf("HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, http.StatusBadGateway)
	}
	if statusCode != http.StatusBadGateway {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadGateway)
	}
	if elapsed >= time.Second {
		t.Errorf("took %v, want well under 1s (the backoff wait should be cut short by ctx)", elapsed)
	}
}

func TestDoRawInternalCtxNeverAnsweringServerHonoursDeadline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block on the request's own context instead of never returning, so
		// srv.Close() below does not hang waiting for this handler to finish.
		<-r.Context().Done()
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, _, err := us.doRawInternalCtx(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, true)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
	}
	if elapsed >= time.Second {
		t.Errorf("took %v, want well under 1s", elapsed)
	}
}

func TestDoRawInternalCtxTokenRefreshCarriesContext(t *testing.T) {
	// The refresh endpoint never answers; it must block on its own request's
	// context (not just never call w.Write) so refreshServer.Close() doesn't hang.
	refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer refreshServer.Close()

	// The main server always answers 401, so doRawInternalCtx always reaches the
	// refresh path.
	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer mainServer.Close()

	// TokenRefresh always dials "https://"+jwt.Domain, so the refresh URL must
	// point at the TLS server, independently of mainHost.
	domain := strings.TrimPrefix(refreshServer.URL, "https://")
	us := New(testJWT(t, domain), "", mainServer.URL)
	us.client = refreshServer.Client() // trust the TLS test server's certificate
	// A safety net, not the mechanism under test: ctx (200ms below) should always win this
	// race. Without it, a regression that stops carrying ctx into the refresh request would
	// hang this test until the whole run's -timeout, instead of failing it in a few seconds.
	us.client.Timeout = 2 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, statusCode, err := us.doRawInternalCtx(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, false)
	elapsed := time.Since(start)

	// ctx ends during the refresh itself, before it can fail or succeed, and no 429/5xx was
	// ever seen on the original request, so this is the done-context case (status 0), not a
	// refresh failure (which would carry the original request's 401).
	if statusCode != 0 {
		t.Errorf("statusCode = %d, want 0", statusCode)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded in its chain (the refresh request should abort with ctx)", err)
	}
	if elapsed >= time.Second {
		t.Errorf("took %v, want well under 1s: the refresh request must carry ctx instead of hanging for the client timeout", elapsed)
	}
}

// --- GetFieldsCtx / GetEntityCtx ---

func TestGetFieldsCtxRequestAndDecode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/crm/v1/entities/leads/fields" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/crm/v1/entities/leads/fields")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer token")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[
			{"code":"source","type":"list","values":[{"value":"ads","title":"Ads","active":false}]},
			{"code":"city","type":"list","dependency":{"id":3,"parent_field_code":"country"}}
		]}`)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	fields, err := us.GetFieldsCtx(context.Background(), "leads")
	if err != nil {
		t.Fatalf("GetFieldsCtx() error = %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("len(fields) = %d, want 2", len(fields))
	}
	if fields[0].Code != "source" || len(fields[0].Values) != 1 {
		t.Fatalf("fields[0] = %+v", fields[0])
	}
	if active := fields[0].Values[0].Active; active == nil || *active {
		t.Errorf("fields[0].Values[0].Active = %v, want a pointer to false", active)
	}
	if dep := fields[1].Dependency; dep == nil || dep.ID != 3 || dep.ParentFieldCode != "country" {
		t.Errorf("fields[1].Dependency = %+v, want {ID:3 ParentFieldCode:country}", dep)
	}
}

func TestGetEntityCtxRawBody(t *testing.T) {
	const respBody = `{"id":5,"first_name":"Олена","phone":[{"value":"380664823817"}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/crm/v1/entities/contacts/5" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/crm/v1/entities/contacts/5")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, respBody)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	body, err := us.GetEntityCtx(context.Background(), "contacts", 5)
	if err != nil {
		t.Fatalf("GetEntityCtx() error = %v", err)
	}
	if string(body) != respBody {
		t.Errorf("body = %q, want %q", body, respBody)
	}
}

// Both methods must normalize a final non-2xx answer to *HTTPError and must not retry a
// 4xx: same non-retry contract TestDoRawRetries5xxButNot4xx pins for the untyped methods.
func TestGetFieldsCtxAndGetEntityCtxHTTPErrorNoRetries(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var requests int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				w.WriteHeader(status)
			}))
			defer server.Close()
			us := New("token", "", server.URL)

			_, err := us.GetFieldsCtx(context.Background(), "leads")
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) {
				t.Fatalf("GetFieldsCtx() error = %v, want *HTTPError", err)
			}
			if httpErr.StatusCode != status {
				t.Errorf("GetFieldsCtx() HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, status)
			}
			if got := atomic.LoadInt32(&requests); got != 1 {
				t.Errorf("GetFieldsCtx(): server received %d requests, want 1 (no retries on %d)", got, status)
			}

			atomic.StoreInt32(&requests, 0)
			_, err = us.GetEntityCtx(context.Background(), "contacts", 5)
			if !errors.As(err, &httpErr) {
				t.Fatalf("GetEntityCtx() error = %v, want *HTTPError", err)
			}
			if httpErr.StatusCode != status {
				t.Errorf("GetEntityCtx() HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, status)
			}
			if got := atomic.LoadInt32(&requests); got != 1 {
				t.Errorf("GetEntityCtx(): server received %d requests, want 1 (no retries on %d)", got, status)
			}
		})
	}
}

// A done context must stop both methods promptly instead of running out the client's
// 30s timeout, the same way it stops doRawInternalCtx directly (see
// TestDoRawInternalCtxNeverAnsweringServerHonoursDeadline).
func TestGetFieldsCtxAndGetEntityCtxDoneContextStops(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block on the request's own context instead of never returning, so
		// srv.Close() below does not hang waiting for this handler to finish.
		<-r.Context().Done()
	}))
	defer server.Close()
	us := New("token", "", server.URL)

	t.Run("GetFieldsCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

func TestGetFieldsCtx401RefreshFailureSurfacesAsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// "token" has no dots, so UnmarshalTokenData rejects it before TokenRefresh ever makes
	// a network call, same setup as TestDoRaw401RefreshFailureReturnsRawError. That raw
	// refresh error must surface through GetFieldsCtx as an *HTTPError, not unwrapped.
	us := New("token", "", server.URL)
	wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, "leads", ""))

	_, err := us.GetFieldsCtx(context.Background(), "leads")

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("GetFieldsCtx() error = %v, want *HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", httpErr.StatusCode, http.StatusUnauthorized)
	}
	if httpErr.Method != http.MethodGet || httpErr.URL != wantURL {
		t.Errorf("Method/URL = %q %q, want %q %q", httpErr.Method, httpErr.URL, http.MethodGet, wantURL)
	}
	wantCauseText := "invalid JWT token format: expected 3 parts, got 1"
	if httpErr.Err == nil || httpErr.Err.Error() != wantCauseText {
		t.Errorf("httpErr.Err = %v, want %q", httpErr.Err, wantCauseText)
	}
	if !errors.Is(err, httpErr.Err) {
		t.Error("errors.Is(err, httpErr.Err) = false, want true: Unwrap should expose the raw refresh error")
	}
}

// A refresh request that never answers must still honour ctx: both methods report the
// context error, not a 401 *HTTPError, once ctx's deadline passes.
func TestGetFieldsCtxAndGetEntityCtxRefreshStallIsContextError(t *testing.T) {
	newClient := func() (*Uspacy, *httptest.Server, *httptest.Server) {
		// The refresh endpoint never answers; it must block on its own request's context
		// (not just never call w.Write) so refreshServer.Close() doesn't hang.
		refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		domain := strings.TrimPrefix(refreshServer.URL, "https://")
		us := New(testJWT(t, domain), "", mainServer.URL)
		us.client = refreshServer.Client() // trust the TLS test server's certificate
		us.client.Timeout = 2 * time.Second
		return us, refreshServer, mainServer
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

// Same shape as TestGetFieldsCtxAndGetEntityCtxRefreshStallIsContextError, but the caller
// cancels ctx instead of letting a deadline pass.
func TestGetFieldsCtxAndGetEntityCtxRefreshCancelIsContextError(t *testing.T) {
	newClient := func() (*Uspacy, *httptest.Server, *httptest.Server) {
		refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		domain := strings.TrimPrefix(refreshServer.URL, "https://")
		us := New(testJWT(t, domain), "", mainServer.URL)
		us.client = refreshServer.Client()
		us.client.Timeout = 2 * time.Second
		return us, refreshServer, mainServer
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(100*time.Millisecond, cancel)

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(100*time.Millisecond, cancel)

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

// A 429 seen before the 401 must still win when the retry that follows a successful-looking
// 401 handoff never completes because the refresh itself stalls past ctx: the done-context
// rule prefers a real HTTP error already on hand over the bare context error.
func TestGetFieldsCtxAndGetEntityCtxRefreshStallKeepsEarlier429(t *testing.T) {
	const respBody = "slow down"
	newClient := func() (*Uspacy, *httptest.Server, *httptest.Server) {
		refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		var calls int32
		mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if atomic.AddInt32(&calls, 1) == 1 {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, respBody)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}))
		domain := strings.TrimPrefix(refreshServer.URL, "https://")
		us := New(testJWT(t, domain), "", mainServer.URL)
		us.client = refreshServer.Client()
		us.client.Timeout = 2 * time.Second
		return us, refreshServer, mainServer
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

// A response whose body never finishes arriving must still honour ctx: if an earlier attempt
// saw a 429, that *HTTPError wins over the bare context error once ctx ends mid-read.
func TestGetFieldsCtxAndGetEntityCtxBodyStallKeepsEarlier429(t *testing.T) {
	const respBody = "slow down"
	newServer := func() *httptest.Server {
		var calls int32
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if atomic.AddInt32(&calls, 1) == 1 {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, respBody)
				return
			}
			// Headers arrive (200), but the body itself never does.
			w.WriteHeader(http.StatusOK)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			<-r.Context().Done()
		}))
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

// A response whose body never finishes arriving, with no earlier 429/5xx to fall back on,
// must report ctx's own error instead of the raw, unwrapped read error with status 0.
func TestGetFieldsCtxAndGetEntityCtxBodyStallIsContextError(t *testing.T) {
	newServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			<-r.Context().Done()
		}))
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		// abortedWhileWaiting's wrapping, not the raw, unwrapped read error: a plain
		// "context deadline exceeded" (io.ReadAll's error passed straight through) would
		// also satisfy errors.Is above, so the prefix is what actually distinguishes the
		// two.
		if !strings.HasPrefix(err.Error(), "request aborted: ") {
			t.Errorf("err = %q, want the \"request aborted: \" prefix", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if !strings.HasPrefix(err.Error(), "request aborted: ") {
			t.Errorf("err = %q, want the \"request aborted: \" prefix", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

// --- W1/W3 fixes: last-HTTP-error tracking, 3xx, connection failures ---

func TestGetFieldsCtx_CancelStopsRetrySleep(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := New("t", "", srv.URL).GetFieldsCtx(ctx, "leads")
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusBadGateway || time.Since(start) > time.Second {
		t.Fatalf("want *HTTPError 502 within 1s (no 3 s backoff), got %v after %v", err, time.Since(start))
	}
	if !strings.HasPrefix(he.Error(), "request failed: [GET] ") || !strings.Contains(he.Error(), "status code: 502, response: ") {
		t.Fatalf("error text changed: %q", he.Error())
	}
}

// Same shape as TestGetFieldsCtx_CancelStopsRetrySleep, for GetEntityCtx. The response body
// is non-empty so this also pins W1(b): the *HTTPError returned when ctx ends must carry the
// last 5xx response's real body, not an empty, synthesized one.
func TestGetEntityCtx_CancelStopsRetrySleep(t *testing.T) {
	const respBody = "bad gateway"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, respBody)
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := New("t", "", srv.URL).GetEntityCtx(ctx, "contacts", 5)
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusBadGateway || time.Since(start) > time.Second {
		t.Fatalf("want *HTTPError 502 within 1s (no 3 s backoff), got %v after %v", err, time.Since(start))
	}
	if string(he.Body) != respBody {
		t.Fatalf("HTTPError.Body = %q, want %q (the real 502 body, not a synthesized empty one)", he.Body, respBody)
	}
	if !strings.HasPrefix(he.Error(), "request failed: [GET] ") || !strings.Contains(he.Error(), "status code: 502, response: "+respBody) {
		t.Fatalf("error text changed: %q", he.Error())
	}
}

// A final 500 must also surface as *HTTPError from both context-aware methods, without
// waiting out the real ~3s backoff.
func TestGetFieldsCtxAndGetEntityCtx500IsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	t.Run("GetFieldsCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("err = %v, want *HTTPError 500", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("err = %v, want *HTTPError 500", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})
}

// A 429 with a long Retry-After must not make the caller actually wait it out: ctx ending
// during that wait still returns the 429 *HTTPError with its body, in well under 1s.
func TestGetFieldsCtxAndGetEntityCtx429WithRetryAfterIsHTTPError(t *testing.T) {
	const respBody = "slow down"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "300")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, respBody)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	t.Run("GetFieldsCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests {
			t.Fatalf("err = %v, want *HTTPError 429", err)
		}
		if string(he.Body) != respBody {
			t.Errorf("HTTPError.Body = %q, want %q", he.Body, respBody)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s (no 300s Retry-After wait)", time.Since(start))
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests {
			t.Fatalf("err = %v, want *HTTPError 429", err)
		}
		if string(he.Body) != respBody {
			t.Errorf("HTTPError.Body = %q, want %q", he.Body, respBody)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s (no 300s Retry-After wait)", time.Since(start))
		}
	})
}

// A closed server (connection refused on every attempt) with a short ctx must stop with the
// context error, not run out the client's 30s timeout or the real ~3s/~5s backoffs.
func TestGetFieldsCtxAndGetEntityCtxConnRefusedIsContextError(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := dead.URL
	dead.Close() // closed: every dial now fails with "connection refused"
	us := New("token", "", deadURL)

	t.Run("GetFieldsCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})
}

// F1: whatever the refresh endpoint answers with (a plain HTTP failure, here), GetFieldsCtx
// and GetEntityCtx must report a 401 *HTTPError for the original request and method/URL, with
// the refresh's own error as its cause — never the refresh's own status or URL directly. The
// old (non-context) path must keep returning that raw refresh error unwrapped and unchanged.
//
// The ctx methods run with a short ctx, not context.Background(): a 4xx from the refresh
// endpoint fails it immediately, well within that budget, so its own *HTTPError becomes the
// 401's cause as above. A 500 is different: doRawInternalCtx retries it like any other 5xx,
// and the short ctx ends during that retry's backoff before the refresh can produce a
// conclusive answer — so ctx ending wins instead, the same done-context rule as everywhere
// else ctx ends mid-retry, with no *HTTPError from the refresh in the chain.
func TestGetFieldsCtxAndGetEntityCtxRefreshHTTPFailureIs401HTTPError(t *testing.T) {
	for _, refreshStatus := range []int{http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError} {
		t.Run(http.StatusText(refreshStatus), func(t *testing.T) {
			newClient := func() (*Uspacy, *httptest.Server) {
				srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, "/auth/refresh_token") {
						w.WriteHeader(refreshStatus)
						fmt.Fprint(w, "refresh rejected")
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				}))
				domain := strings.TrimPrefix(srv.URL, "https://")
				us := New(testJWT(t, domain), "", srv.URL)
				us.client = srv.Client()
				return us, srv
			}

			t.Run("GetFieldsCtx", func(t *testing.T) {
				us, srv := newClient()
				defer srv.Close()
				wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, "leads", ""))
				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
				defer cancel()

				_, err := us.GetFieldsCtx(ctx, "leads")
				if refreshStatus == http.StatusInternalServerError {
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
					}
					var he *HTTPError
					if errors.As(err, &he) {
						t.Errorf("err = %v, want no *HTTPError in the chain", err)
					}
					return
				}
				var he *HTTPError
				if !errors.As(err, &he) {
					t.Fatalf("err = %v, want *HTTPError", err)
				}
				if he.Method != http.MethodGet || he.URL != wantURL || he.StatusCode != http.StatusUnauthorized {
					t.Errorf("HTTPError = %+v, want Method=GET URL=%q StatusCode=401", he, wantURL)
				}
				var cause *HTTPError
				if !errors.As(he.Err, &cause) || cause.StatusCode != refreshStatus || cause.Method != http.MethodPost {
					t.Errorf("HTTPError.Err = %v, want the refresh's own *HTTPError with Method=POST StatusCode=%d", he.Err, refreshStatus)
				}
			})

			t.Run("GetEntityCtx", func(t *testing.T) {
				us, srv := newClient()
				defer srv.Close()
				wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, "contacts"), "5")
				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
				defer cancel()

				_, err := us.GetEntityCtx(ctx, "contacts", 5)
				if refreshStatus == http.StatusInternalServerError {
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
					}
					var he *HTTPError
					if errors.As(err, &he) {
						t.Errorf("err = %v, want no *HTTPError in the chain", err)
					}
					return
				}
				var he *HTTPError
				if !errors.As(err, &he) {
					t.Fatalf("err = %v, want *HTTPError", err)
				}
				if he.Method != http.MethodGet || he.URL != wantURL || he.StatusCode != http.StatusUnauthorized {
					t.Errorf("HTTPError = %+v, want Method=GET URL=%q StatusCode=401", he, wantURL)
				}
				var cause *HTTPError
				if !errors.As(he.Err, &cause) || cause.StatusCode != refreshStatus || cause.Method != http.MethodPost {
					t.Errorf("HTTPError.Err = %v, want the refresh's own *HTTPError with Method=POST StatusCode=%d", he.Err, refreshStatus)
				}
			})

			// The old path has no ctx to cut a retried 500 short, and that combination
			// (a 5xx retried with real backoff under context.Background()) is already
			// pinned by TestDoRawRetries5xxButNot4xx and, for the refresh-failure-stays-
			// raw invariant itself, by TestDoRaw401RefreshFailureReturnsRawError. Skip it
			// here rather than pay a real ~8s backoff for a case already covered.
			if refreshStatus == http.StatusInternalServerError {
				return
			}
			t.Run("GetFields_OldPathUnwrapped", func(t *testing.T) {
				us, srv := newClient()
				defer srv.Close()

				_, err := us.GetFields("leads")
				var he *HTTPError
				if !errors.As(err, &he) {
					t.Fatalf("err = %v, want *HTTPError (the refresh's own, unwrapped)", err)
				}
				if he.StatusCode != refreshStatus || he.Method != http.MethodPost {
					t.Errorf("GetFields() err = %+v, want the raw refresh error: Method=POST StatusCode=%d", he, refreshStatus)
				}
				if he.StatusCode == http.StatusUnauthorized {
					t.Errorf("GetFields() err has StatusCode=401: the old path must not gain the ctx path's 401 wrapping")
				}
			})
		})
	}
}

// W1(a): a stale pre-refresh 401 must not resurface as the final answer when the retried
// request (after a successful refresh) then fails in flight because ctx ends. The result must
// be a context error, with no *HTTPError in its chain.
func TestGetFieldsCtxAndGetEntityCtxRefreshedRetryAbortedInFlight(t *testing.T) {
	newServer := func() (*httptest.Server, string) {
		var mainCalls int32
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/auth/refresh_token") {
				fmt.Fprint(w, `{"jwt":"a.b.c","refreshToken":"r"}`)
				return
			}
			if atomic.AddInt32(&mainCalls, 1) == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			<-r.Context().Done() // the retried request: hang until ctx ends
		}))
		return srv, strings.TrimPrefix(srv.URL, "https://")
	}

	t.Run("GetFieldsCtx", func(t *testing.T) {
		srv, domain := newServer()
		defer srv.Close()
		us := New(testJWT(t, domain), "", srv.URL)
		us.client = srv.Client()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFieldsCtx(ctx, "leads")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain (the stale pre-refresh 401 must not resurface)", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})

	t.Run("GetEntityCtx", func(t *testing.T) {
		srv, domain := newServer()
		defer srv.Close()
		us := New(testJWT(t, domain), "", srv.URL)
		us.client = srv.Client()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntityCtx(ctx, "contacts", 5)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		var he *HTTPError
		if errors.As(err, &he) {
			t.Errorf("err = %v, want no *HTTPError in the chain (the stale pre-refresh 401 must not resurface)", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})
}

// W1(c): ctx ending during the LAST attempt's request — which runs no backoff sleep
// afterward — must still report the last 429/5xx seen, not fall through to a stale status
// left by a still-earlier attempt. Retry-After: 0 advances the first two (of defaultRetries=3)
// attempts in about 0s so the test stays fast; the third and last attempt hangs until ctx
// ends.
func TestGetEntityCtxLastAttemptInFlightAbortAfterEarlier429(t *testing.T) {
	var calls int32
	const respBody = "slow down"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, respBody)
			return
		}
		<-r.Context().Done()
	}))
	defer srv.Close()

	us := New("token", "", srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := us.GetEntityCtx(ctx, "contacts", 5)
	elapsed := time.Since(start)

	var he *HTTPError
	if !errors.As(err, &he) {
		t.Fatalf("err = %v, want *HTTPError", err)
	}
	if he.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want %d", he.StatusCode, http.StatusTooManyRequests)
	}
	if string(he.Body) != respBody {
		t.Errorf("Body = %q, want %q", he.Body, respBody)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("server received %d requests, want 3 (all three attempts used)", got)
	}
	if elapsed > time.Second {
		t.Errorf("took %v, want well under 1s", elapsed)
	}
}

// Two 429s, then a refreshed 401 right before the last attempt, which is cut off in flight.
// The last real 429 must win over both the stale 401 and the bare ctx error: lastHTTPErr must
// track the outer loop's own attempts (not the 401 a successful refresh resolves), and the
// last attempt's own ctx check must still fire after a refresh mid-loop.
func TestGetEntityCtxLastAttemptInFlightAfterRefreshed401(t *testing.T) {
	var calls int32
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/auth/refresh_token") {
			fmt.Fprint(w, `{"jwt":"a.b.c","refreshToken":"r"}`)
			return
		}
		switch n := atomic.AddInt32(&calls, 1); {
		case n <= 2:
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, "slow down")
		case n == 3:
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "stale")
		default:
			<-r.Context().Done()
		}
	}))
	defer srv.Close()
	us := New(testJWT(t, strings.TrimPrefix(srv.URL, "https://")), "", srv.URL)
	us.client = srv.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	_, err := us.GetEntityCtx(ctx, "contacts", 5)
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != "slow down" {
		t.Fatalf("err = %v, want the last 429 *HTTPError with its body", err)
	}
}

// A 502, then a later, non-last attempt is cut off in flight (ctx outlasts the first backoff).
// The real 502 must come back, body included, rather than the bare ctx error the sleep itself
// returned.
func TestGetEntityCtx502ThenLaterAttemptInFlight(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "bad gateway")
			return
		}
		<-r.Context().Done()
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	_, err := New("token", "", srv.URL).GetEntityCtx(ctx, "contacts", 5)
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusBadGateway || string(he.Body) != "bad gateway" {
		t.Fatalf("err = %v, want the 502 *HTTPError with its body", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("server received %d requests, want 2", got)
	}
}

// A final 3xx (e.g. a 302 with no Location) is success for every non-context method, but
// GetFieldsCtx and GetEntityCtx must report it as an *HTTPError.
func TestGetFieldsCtxAndGetEntityCtx302IsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound) // no Location header
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	_, err := us.GetFieldsCtx(context.Background(), "leads")
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusFound {
		t.Errorf("GetFieldsCtx() err = %v, want *HTTPError 302", err)
	}

	_, err = us.GetEntityCtx(context.Background(), "contacts", 5)
	if !errors.As(err, &he) || he.StatusCode != http.StatusFound {
		t.Errorf("GetEntityCtx() err = %v, want *HTTPError 302", err)
	}
}

// The mirror of TestGetFieldsCtxAndGetEntityCtx302IsHTTPError: a final 3xx stays a success on
// the non-context methods, both through doRaw directly and through a method built on it.
func TestDoRawOldPath302IsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound) // no Location header
		fmt.Fprint(w, "moved")
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)
	body, statusCode, err := us.doRaw(us.buildURL("resource"), http.MethodGet, nil, nil)
	if err != nil || statusCode != http.StatusFound || string(body) != "moved" {
		t.Errorf("doRaw() = %q, %d, %v; want \"moved\", 302, nil", body, statusCode, err)
	}
	if b, err := us.GetList("contacts", nil); err != nil || string(b) != "moved" {
		t.Errorf("GetList() = %q, %v; want \"moved\", nil", b, err)
	}
}

// Characterization: a 429 with Retry-After: 0 lets all defaultRetries attempts run almost
// immediately; pins the exact final error text and request count for the non-context path,
// which the context-aware retry/backoff changes must not touch.
func TestDoRaw429RetryAfterZeroFinalErrorText(t *testing.T) {
	var requests int32
	const respBody = "slow down"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, respBody)
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	url := us.buildURL("resource")
	body, statusCode, err := us.doRaw(url, http.MethodGet, nil, nil)

	if statusCode != http.StatusTooManyRequests {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusTooManyRequests)
	}
	if string(body) != respBody {
		t.Errorf("body = %q, want %q", body, respBody)
	}
	wantErr := fmt.Sprintf("request failed: [%s] %s, status code: %d, response: %s", http.MethodGet, url, http.StatusTooManyRequests, respBody)
	if err == nil || err.Error() != wantErr {
		t.Errorf("err = %v, want %q", err, wantErr)
	}
	if got := atomic.LoadInt32(&requests); got != defaultRetries {
		t.Errorf("server received %d requests, want %d", got, defaultRetries)
	}
}

// Regression for the lastStatusCode field removed from Uspacy: GetFieldsCtx and GetEntityCtx
// running in parallel on one *Uspacy — the documented usage (a portal client adapter shared
// by both reads) — must not race on shared state. Meaningful under go test -race.
func TestGetFieldsCtxAndGetEntityCtxParallelNoRace(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, _ = us.GetFieldsCtx(context.Background(), "leads")
	}()
	go func() {
		defer wg.Done()
		_, _ = us.GetEntityCtx(context.Background(), "contacts", 5)
	}()
	wg.Wait()
}
