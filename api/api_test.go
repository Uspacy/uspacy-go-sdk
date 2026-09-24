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

func TestAsHTTPError(t *testing.T) {
	const method, url = http.MethodGet, "https://example.uspacy.ua/x"

	if got := asHTTPError(method, url, http.StatusOK, nil); got != nil {
		t.Errorf("asHTTPError(200, nil) = %v, want nil", got)
	}

	raw := errors.New("boom")
	got := asHTTPError(method, url, http.StatusUnauthorized, raw)
	var httpErr *HTTPError
	if !errors.As(got, &httpErr) {
		t.Fatalf("asHTTPError(401, raw) = %v, want *HTTPError", got)
	}
	if httpErr.Method != method || httpErr.URL != url || httpErr.StatusCode != http.StatusUnauthorized || httpErr.Err != raw {
		t.Errorf("asHTTPError(401, raw) = %+v, want Method=%q URL=%q StatusCode=401 Err=raw", httpErr, method, url)
	}

	already := &HTTPError{StatusCode: http.StatusInternalServerError}
	if got := asHTTPError(method, url, http.StatusInternalServerError, already); got != already {
		t.Error("asHTTPError should return an existing *HTTPError unchanged, not re-wrap it")
	}

	if got := asHTTPError(method, url, http.StatusOK, raw); got != raw {
		t.Errorf("asHTTPError(200, raw) = %v, want raw unchanged (status < 400)", got)
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

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, statusCode, err := us.doRawInternalCtx(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, false)
	elapsed := time.Since(start)

	if statusCode != http.StatusUnauthorized {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusUnauthorized)
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
