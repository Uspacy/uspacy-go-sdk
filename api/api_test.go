package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
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
	productList, err := us.GetEntityProductList(context.Background(), "deals", 1072)
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
	product, err := us.CreateEntityListProduct(context.Background(), crm.CreateEntityListProductRequest{
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
	statusCode, err := us.DeleteEntityListProduct(context.Background(), 9)
	if err != nil {
		t.Fatalf("DeleteEntityListProduct() error = %v", err)
	}
	if statusCode != http.StatusNoContent {
		t.Errorf("DeleteEntityListProduct() status = %d, want %d", statusCode, http.StatusNoContent)
	}
}

// --- Characterization tests for doRequest/doRequest ---
func TestDoRawRetries5xxButNot4xx(t *testing.T) {
	cases := []struct {
		name         string
		status       int
		wantRequests int32
	}{
		{"5xx is retried up to defaultMaxRetries", http.StatusInternalServerError, defaultMaxRetries},
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

			us := New("token", "", server.URL, WithRetryBackoff(time.Millisecond, 5*time.Millisecond))
			_, statusCode, err := us.doRaw(context.Background(), us.buildURL("resource"), http.MethodGet, nil, nil)

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

func TestDoRaw401RefreshFailureIsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// "token" has no dots, so UnmarshalTokenData rejects it before TokenRefresh ever
	// makes a network call. The 401 status and original request details are preserved
	// as an *HTTPError; the parse error is attached via Err.
	us := New("token", "", server.URL)
	_, statusCode, err := us.doRaw(context.Background(), us.buildURL("resource"), http.MethodGet, nil, nil)

	if statusCode != http.StatusUnauthorized {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusUnauthorized)
	}
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusUnauthorized {
		t.Fatalf("err = %v, want *HTTPError 401", err)
	}
	wantErr := "invalid JWT token format: expected 3 parts, got 1"
	if he.Err == nil || he.Err.Error() != wantErr {
		t.Errorf("err.Err = %v, want %q", he.Err, wantErr)
	}
}

// A refresh endpoint failure after a 401 is reported as a 401 *HTTPError on the original
// request, with the refresh error as its cause.
func TestDoRaw401RefreshHTTPFailureIsHTTPError(t *testing.T) {
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

	_, _, err := us.doRaw(context.Background(), us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, "leads", "")), http.MethodGet, nil, nil)
	he, ok := err.(*HTTPError)
	if !ok || he.StatusCode != http.StatusUnauthorized || he.Method != http.MethodGet {
		t.Errorf("doRaw() err = %T %v, want *HTTPError 401 for the original request", err, err)
	}
	var cause *HTTPError
	if !errors.As(he.Err, &cause) || cause.StatusCode != http.StatusForbidden || cause.Method != http.MethodPost {
		t.Errorf("err.Err = %v, want the refresh's own *HTTPError (POST 403)", he.Err)
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
	body, statusCode, err := us.doRaw(context.Background(), url, http.MethodGet, nil, nil)

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

// requestFailedAfterRetries builds the error for doRequest's post-loop path, taken
// when every attempt fails at the transport level (statusCode stays 0). Exercising that
// exact path through doRequest itself would need two real backoff sleeps
// (defaultMaxRetries=3, with jittered ~3s/~5s delays) to run to completion without ctx firing,
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
		err := requestFailedAfterRetries(context.Background(), defaultMaxRetries, logs)
		wantText := fmt.Sprintf("request failed after %d retries:\n%s", defaultMaxRetries, wantDetails)
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

		err := requestFailedAfterRetries(ctx, defaultMaxRetries, logs)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled in its chain", err)
		}
		wantPrefix := fmt.Sprintf("request failed after %d retries:\n%scontext error: ", defaultMaxRetries, wantDetails)
		if !strings.HasPrefix(err.Error(), wantPrefix) {
			t.Errorf("Error() = %q, want prefix %q", err.Error(), wantPrefix)
		}
	})

	t.Run("expired deadline is reachable via errors.Is", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()
		<-ctx.Done()

		err := requestFailedAfterRetries(ctx, defaultMaxRetries, logs)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
	})
}

// --- Context-path tests for doRequest ---

func TestSleep(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	// An already-done ctx wins at every wait, including zero and negative ones (Retry-After: 0
	// or a negative value). Without the up-front check, select could pick an already-fired
	// timer over ctx.Done(), so repeat the call enough times to catch that.
	for _, d := range []time.Duration{0, -time.Second, time.Hour} {
		for i := 0; i < 200; i++ {
			if err := sleep(canceled, d); !errors.Is(err, context.Canceled) {
				t.Fatalf("sleep(canceled, %v) = %v, want context.Canceled", d, err)
			}
		}
	}

	// A live ctx with a non-positive wait returns nil at once.
	for _, d := range []time.Duration{0, -time.Second} {
		start := time.Now()
		if err := sleep(context.Background(), d); err != nil || time.Since(start) > 50*time.Millisecond {
			t.Fatalf("sleep(background, %v) = %v after %v, want nil at once", d, err, time.Since(start))
		}
	}

	// A live ctx waits the full d.
	start := time.Now()
	if err := sleep(context.Background(), 20*time.Millisecond); err != nil || time.Since(start) < 20*time.Millisecond {
		t.Fatalf("sleep(background, 20ms) = %v after %v, want nil after at least 20ms", err, time.Since(start))
	}

	// A ctx that ends mid-wait stops the wait with its own error.
	short, cancelShort := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelShort()
	start = time.Now()
	if err := sleep(short, time.Hour); !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("sleep(deadline, 1h) = %v after %v, want context.DeadlineExceeded within 1s", err, time.Since(start))
	}
}

// testJWT builds a syntactically valid but unsigned JWT whose payload sets the
// "domain" claim, which is all UnmarshalTokenData/tokenRefresh need from it.
func testJWT(t *testing.T, domain string) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"domain":%q}`, domain)))
	return header + "." + payload + ".sig"
}

func TestDoRawInternalRetryWaitCutShortBy502(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "bad gateway")
	}))
	defer server.Close()

	us := New("token", "", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, statusCode, err := us.doRequest(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, true)
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

func TestDoRawInternalNeverAnsweringServerHonoursDeadline(t *testing.T) {
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
	_, _, err := us.doRequest(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, true)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
	}
	if elapsed >= time.Second {
		t.Errorf("took %v, want well under 1s", elapsed)
	}
}

func TestDoRawInternalTokenRefreshCarriesContext(t *testing.T) {
	// The refresh endpoint never answers; it must block on its own request's
	// context (not just never call w.Write) so refreshServer.Close() doesn't hang.
	refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer refreshServer.Close()
	defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout

	// The main server always answers 401, so doRequest always reaches the
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
	_, statusCode, err := us.doRequest(ctx, us.buildURL("resource"), http.MethodGet, nil, nil, false)
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

// --- GetFields / GetEntity ---

func TestGetFieldsRequestAndDecode(t *testing.T) {
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
	fields, err := us.GetFields(context.Background(), "leads")
	if err != nil {
		t.Fatalf("GetFields() error = %v", err)
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

func TestGetEntityRawBody(t *testing.T) {
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
	body, err := us.GetEntity(context.Background(), "contacts", 5)
	if err != nil {
		t.Fatalf("GetEntity() error = %v", err)
	}
	if string(body) != respBody {
		t.Errorf("body = %q, want %q", body, respBody)
	}
}

// Both methods must normalize a final non-2xx answer to *HTTPError and must not retry a
// 4xx: same non-retry contract TestDoRawRetries5xxButNot4xx pins for the untyped methods.
func TestGetFieldsAndGetEntityHTTPErrorNoRetries(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var requests int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				w.WriteHeader(status)
			}))
			defer server.Close()
			us := New("token", "", server.URL)

			_, err := us.GetFields(context.Background(), "leads")
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) {
				t.Fatalf("GetFields() error = %v, want *HTTPError", err)
			}
			if httpErr.StatusCode != status {
				t.Errorf("GetFields() HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, status)
			}
			if got := atomic.LoadInt32(&requests); got != 1 {
				t.Errorf("GetFields(): server received %d requests, want 1 (no retries on %d)", got, status)
			}

			atomic.StoreInt32(&requests, 0)
			_, err = us.GetEntity(context.Background(), "contacts", 5)
			if !errors.As(err, &httpErr) {
				t.Fatalf("GetEntity() error = %v, want *HTTPError", err)
			}
			if httpErr.StatusCode != status {
				t.Errorf("GetEntity() HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, status)
			}
			if got := atomic.LoadInt32(&requests); got != 1 {
				t.Errorf("GetEntity(): server received %d requests, want 1 (no retries on %d)", got, status)
			}
		})
	}
}

// A done context must stop both methods promptly instead of running out the client's
// 30s timeout, the same way it stops doRequest directly (see
// TestDoRawInternalNeverAnsweringServerHonoursDeadline).
func TestGetFieldsAndGetEntityDoneContextStops(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block on the request's own context instead of never returning, so
		// srv.Close() below does not hang waiting for this handler to finish.
		<-r.Context().Done()
	}))
	defer server.Close()
	us := New("token", "", server.URL)

	t.Run("GetFields", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntity", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
		elapsed := time.Since(start)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context.DeadlineExceeded in its chain", err)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})
}

func TestGetFields401RefreshFailureSurfacesAsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// "token" has no dots, so UnmarshalTokenData rejects it before TokenRefresh ever makes
	// a network call, same setup as TestDoRaw401RefreshFailureReturnsRawError. That raw
	// refresh error must surface through GetFields as an *HTTPError, not unwrapped.
	us := New("token", "", server.URL)
	wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, "leads", ""))

	_, err := us.GetFields(context.Background(), "leads")

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("GetFields() error = %v, want *HTTPError", err)
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
func TestGetFieldsAndGetEntityRefreshStallIsContextError(t *testing.T) {
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

	t.Run("GetFields", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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

// Same shape as TestGetFieldsAndGetEntityRefreshStallIsContextError, but the caller
// cancels ctx instead of letting a deadline pass.
func TestGetFieldsAndGetEntityRefreshCancelIsContextError(t *testing.T) {
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

	t.Run("GetFields", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(100*time.Millisecond, cancel)

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(100*time.Millisecond, cancel)

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntityRefreshStallKeepsEarlier429(t *testing.T) {
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

	t.Run("GetFields", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntity", func(t *testing.T) {
		us, refreshServer, mainServer := newClient()
		defer refreshServer.Close()
		defer refreshServer.CloseClientConnections() // end the detached shared refresh now, not at the client timeout
		defer mainServer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntityBodyStallKeepsEarlier429(t *testing.T) {
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

	t.Run("GetFields", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
		elapsed := time.Since(start)

		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != respBody {
			t.Fatalf("err = %v, want the 429 *HTTPError with body %q", err, respBody)
		}
		if elapsed >= time.Second {
			t.Errorf("took %v, want well under 1s", elapsed)
		}
	})

	t.Run("GetEntity", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntityBodyStallIsContextError(t *testing.T) {
	newServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			<-r.Context().Done()
		}))
	}

	t.Run("GetFields", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		srv := newServer()
		defer srv.Close()
		us := New("token", "", srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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

// --- Last-HTTP-error tracking on ctx ending, 3xx, and connection failures ---

func TestGetFields_CancelStopsRetrySleep(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := New("t", "", srv.URL).GetFields(ctx, "leads")
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusBadGateway || time.Since(start) > time.Second {
		t.Fatalf("want *HTTPError 502 within 1s (no 3 s backoff), got %v after %v", err, time.Since(start))
	}
	if !strings.HasPrefix(he.Error(), "request failed: [GET] ") || !strings.Contains(he.Error(), "status code: 502, response: ") {
		t.Fatalf("error text changed: %q", he.Error())
	}
}

// Same shape as TestGetFields_CancelStopsRetrySleep, for GetEntity. The response body
// is non-empty so this also pins that the *HTTPError returned when ctx ends must carry the
// last 5xx response's real body, not an empty, synthesized one.
func TestGetEntity_CancelStopsRetrySleep(t *testing.T) {
	const respBody = "bad gateway"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, respBody)
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := New("t", "", srv.URL).GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntity500IsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	t.Run("GetFields", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
		var he *HTTPError
		if !errors.As(err, &he) || he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("err = %v, want *HTTPError 500", err)
		}
		if time.Since(start) > time.Second {
			t.Errorf("took %v, want well under 1s", time.Since(start))
		}
	})

	t.Run("GetEntity", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntity429WithRetryAfterIsHTTPError(t *testing.T) {
	const respBody = "slow down"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "300")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, respBody)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	t.Run("GetFields", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetFieldsAndGetEntityConnRefusedIsContextError(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := dead.URL
	dead.Close() // closed: every dial now fails with "connection refused"
	us := New("token", "", deadURL)

	t.Run("GetFields", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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

// Whatever the refresh endpoint answers with (a plain HTTP failure, here), GetFields
// and GetEntity must report a 401 *HTTPError for the original request and method/URL, with
// the refresh's own error as its cause — never the refresh's own status or URL directly. The
// A refresh endpoint failure after a 401 is reported as a 401 *HTTPError on the original
// request, with the refresh response as the cause, regardless of the refresh status code.
func TestGetFieldsAndGetEntityRefreshHTTPFailureIs401HTTPError(t *testing.T) {
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

			t.Run("GetFields", func(t *testing.T) {
				us, srv := newClient()
				defer srv.Close()
				wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, "leads", ""))
				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
				defer cancel()

				_, err := us.GetFields(ctx, "leads")
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

			t.Run("GetEntity", func(t *testing.T) {
				us, srv := newClient()
				defer srv.Close()
				wantURL := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, "contacts"), "5")
				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
				defer cancel()

				_, err := us.GetEntity(ctx, "contacts", 5)
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

		})
	}
}

// A stale pre-refresh 401 must not resurface as the final answer when the retried request
// (after a successful refresh) then fails in flight because ctx ends. The result must be a
// context error, with no *HTTPError in its chain.
func TestGetFieldsAndGetEntityRefreshedRetryAbortedInFlight(t *testing.T) {
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

	t.Run("GetFields", func(t *testing.T) {
		srv, domain := newServer()
		defer srv.Close()
		us := New(testJWT(t, domain), "", srv.URL)
		us.client = srv.Client()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetFields(ctx, "leads")
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

	t.Run("GetEntity", func(t *testing.T) {
		srv, domain := newServer()
		defer srv.Close()
		us := New(testJWT(t, domain), "", srv.URL)
		us.client = srv.Client()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		start := time.Now()
		_, err := us.GetEntity(ctx, "contacts", 5)
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

// ctx ending during the LAST attempt's request — which runs no backoff sleep afterward — must
// still report the last 429/5xx seen, not fall through to a stale status left by a
// still-earlier attempt. Retry-After: 0 advances the first two (of defaultMaxRetries=3) attempts
// in about 0s so the test stays fast; the third and last attempt hangs until ctx ends.
func TestGetEntityLastAttemptInFlightAbortAfterEarlier429(t *testing.T) {
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
	_, err := us.GetEntity(ctx, "contacts", 5)
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
func TestGetEntityLastAttemptInFlightAfterRefreshed401(t *testing.T) {
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

	_, err := us.GetEntity(ctx, "contacts", 5)
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusTooManyRequests || string(he.Body) != "slow down" {
		t.Fatalf("err = %v, want the last 429 *HTTPError with its body", err)
	}
}

// A 502, then a later, non-last attempt is cut off in flight (ctx outlasts the first backoff).
// The real 502 must come back, body included, rather than the bare ctx error the sleep itself
// returned.
func TestGetEntity502ThenLaterAttemptInFlight(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	_, err := New("token", "", srv.URL, WithRetryBackoff(10*time.Millisecond, 50*time.Millisecond)).GetEntity(ctx, "contacts", 5)
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusBadGateway || string(he.Body) != "bad gateway" {
		t.Fatalf("err = %v, want the 502 *HTTPError with its body", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("server received %d requests, want 2", got)
	}
}

// A final 3xx is reported as an *HTTPError by every method.
func Test3xxIsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound) // no Location header
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	_, err := us.GetFields(context.Background(), "leads")
	var he *HTTPError
	if !errors.As(err, &he) || he.StatusCode != http.StatusFound {
		t.Errorf("GetFields() err = %v, want *HTTPError 302", err)
	}

	_, err = us.GetEntity(context.Background(), "contacts", 5)
	if !errors.As(err, &he) || he.StatusCode != http.StatusFound {
		t.Errorf("GetEntity() err = %v, want *HTTPError 302", err)
	}

	_, _, err = us.doRaw(context.Background(), us.buildURL("resource"), http.MethodGet, nil, nil)
	if !errors.As(err, &he) || he.StatusCode != http.StatusFound {
		t.Errorf("doRaw() err = %v, want *HTTPError 302", err)
	}
}

// Characterization: a 429 with Retry-After: 0 lets all defaultMaxRetries attempts run almost
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
	body, statusCode, err := us.doRaw(context.Background(), url, http.MethodGet, nil, nil)

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
	if got := atomic.LoadInt32(&requests); got != defaultMaxRetries {
		t.Errorf("server received %d requests, want %d", got, defaultMaxRetries)
	}
}

// Regression for the lastStatusCode field removed from Uspacy: GetFields and GetEntity
// running in parallel on one shared *Uspacy must not race on shared state. Meaningful under
// go test -race.
func TestGetFieldsAndGetEntityParallelNoRace(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer srv.Close()
	us := New("token", "", srv.URL)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, _ = us.GetFields(context.Background(), "leads")
	}()
	go func() {
		defer wg.Done()
		_, _ = us.GetEntity(context.Background(), "contacts", 5)
	}()
	wg.Wait()
}

// --- request options, transport retries, shared token refresh ---

// countingTransport counts round trips, including ones that never reach a server.
type countingTransport struct {
	n  atomic.Int32
	rt http.RoundTripper
}

func (c *countingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.n.Add(1)
	return c.rt.RoundTrip(r)
}

func TestRequestOptionsKeepContentTypeAndReplaceHeaders(t *testing.T) {
	type seen struct {
		contentType string
		auth        []string
		trace       string
	}
	got := make(chan seen, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- seen{r.Header.Get("Content-Type"), r.Header.Values("Authorization"), r.Header.Get("X-Trace")}
		fmt.Fprint(w, `{"id":1}`)
	}))
	defer server.Close()
	us := New("token", "", server.URL)

	if _, _, err := us.CreateEntity(context.Background(), "leads", map[string]any{"x": 1},
		WithHeader("X-Trace", "1"), WithHeader("Authorization", "Bearer other")); err != nil {
		t.Fatalf("CreateEntity() error = %v", err)
	}
	s := <-got
	if s.contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json even with request options", s.contentType)
	}
	if len(s.auth) != 1 || s.auth[0] != "Bearer other" {
		t.Errorf("Authorization = %q, want a single overridden value", s.auth)
	}
	if s.trace != "1" {
		t.Errorf("X-Trace = %q, want 1", s.trace)
	}

	if _, err := us.CreateFile(context.Background(), "leads", "1",
		map[string]io.ReadCloser{"a.txt": io.NopCloser(strings.NewReader("x"))}, WithHeader("X-Trace", "2")); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	s = <-got
	if !strings.HasPrefix(s.contentType, "multipart/form-data") || s.trace != "2" {
		t.Errorf("CreateFile headers = %+v, want multipart Content-Type and X-Trace 2", s)
	}
}

// A POST that timed out after it was sent may already have been applied, so it must not
// be sent again; an idempotent GET still is.
func TestTransportTimeoutRetriesOnlyIdempotent(t *testing.T) {
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		select {
		case <-time.After(time.Second):
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	us := New("token", "", server.URL,
		WithHTTPClient(&http.Client{Timeout: 50 * time.Millisecond}),
		WithRetryBackoff(time.Millisecond, time.Millisecond))

	if _, _, err := us.CreateEntity(context.Background(), "leads", map[string]any{"x": 1}); err == nil {
		t.Fatal("CreateEntity() error = nil, want a timeout")
	}
	if got := hits.Load(); got != 1 {
		t.Errorf("POST reached the server %d times, want 1", got)
	}

	hits.Store(0)
	if _, err := us.GetEntity(context.Background(), "leads", 1); err == nil {
		t.Fatal("GetEntity() error = nil, want a timeout")
	}
	if got := hits.Load(); got != defaultMaxRetries {
		t.Errorf("GET reached the server %d times, want %d", got, defaultMaxRetries)
	}
}

// A POST whose connection was never established cannot have been applied, so it is retried.
func TestTransportDialErrorRetriesPost(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	addr := server.URL
	server.Close()

	ct := &countingTransport{rt: http.DefaultTransport}
	us := New("token", "", addr, WithHTTPClient(&http.Client{Transport: ct}), WithRetryBackoff(time.Millisecond, time.Millisecond))
	if _, _, err := us.CreateEntity(context.Background(), "leads", map[string]any{"x": 1}); err == nil {
		t.Fatal("CreateEntity() error = nil, want connection refused")
	}
	if got := ct.n.Load(); got != defaultMaxRetries {
		t.Errorf("round trips = %d, want %d", got, defaultMaxRetries)
	}
}

// One caller giving up must not fail the shared refresh for another caller still waiting.
func TestSharedRefreshSurvivesFirstCallerCancel(t *testing.T) {
	var refreshHits atomic.Int32
	refreshStarted := make(chan struct{})
	release := make(chan struct{})
	var newJWT string
	refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if refreshHits.Add(1) == 1 {
			close(refreshStarted)
		}
		<-release
		fmt.Fprintf(w, `{"jwt":%q,"refreshToken":"r2"}`, newJWT)
	}))
	defer refreshServer.Close()
	domain := strings.TrimPrefix(refreshServer.URL, "https://")
	newJWT = testJWT(t, domain+".new")

	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != tokenPrefix+newJWT {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer mainServer.Close()

	us := New(testJWT(t, domain), "", mainServer.URL)
	us.client = refreshServer.Client()

	ctxA, cancelA := context.WithCancel(context.Background())
	errA := make(chan error, 1)
	go func() { _, err := us.GetFields(ctxA, "leads"); errA <- err }()
	<-refreshStarted

	errB := make(chan error, 1)
	go func() { _, err := us.GetFields(context.Background(), "leads"); errB <- err }()
	time.Sleep(100 * time.Millisecond) // let B join the refresh in flight
	cancelA()
	if err := <-errA; !errors.Is(err, context.Canceled) {
		t.Errorf("caller A error = %v, want context.Canceled", err)
	}
	close(release)

	if err := <-errB; err != nil {
		t.Errorf("caller B error = %v, want nil: A's cancel must not fail the shared refresh", err)
	}
	if got := refreshHits.Load(); got != 1 {
		t.Errorf("refresh requests = %d, want 1", got)
	}
	if _, refresh := us.Tokens(); refresh != "r2" {
		t.Errorf("Tokens() refresh = %q, want r2", refresh)
	}
}

// The refresh endpoint must get the refresh token, not the (possibly expired) access token,
// and both new tokens must be visible through Tokens.
func TestRefreshSendsRefreshTokenAndUpdatesTokens(t *testing.T) {
	refreshAuth := make(chan string, 1)
	var newJWT string
	refreshServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshAuth <- r.Header.Get("Authorization")
		fmt.Fprintf(w, `{"jwt":%q,"refreshToken":"refresh-2"}`, newJWT)
	}))
	defer refreshServer.Close()
	domain := strings.TrimPrefix(refreshServer.URL, "https://")
	newJWT = testJWT(t, domain+".new")

	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != tokenPrefix+newJWT {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer mainServer.Close()

	us := New(testJWT(t, domain), "Bearer refresh-1", mainServer.URL)
	us.client = refreshServer.Client()
	if _, err := us.GetFields(context.Background(), "leads"); err != nil {
		t.Fatalf("GetFields() error = %v", err)
	}
	if got := <-refreshAuth; got != "Bearer refresh-1" {
		t.Errorf("refresh Authorization = %q, want the refresh token", got)
	}
	if access, refresh := us.Tokens(); access != newJWT || refresh != "refresh-2" {
		t.Errorf("Tokens() = %q, %q, want the refreshed pair", access, refresh)
	}
}

// A 401 for a token that has already been replaced retries with the new token instead of
// refreshing again. "token" is not a JWT, so any refresh attempt would fail.
func TestTokenRefreshSkipsWhenTokenAlreadyChanged(t *testing.T) {
	us := New("token", "", "http://unused.invalid")
	got, err := us.tokenRefresh(context.Background(), "older-token")
	if err != nil || got != "token" {
		t.Fatalf("tokenRefresh(stale) = %q, %v, want the current token and no refresh", got, err)
	}
	if _, err := us.TokenRefresh(context.Background()); err == nil {
		t.Error("TokenRefresh() error = nil, want a forced refresh attempt that fails on the non-JWT token")
	}
}

func TestNextBackoffNoOverflow(t *testing.T) {
	cases := []struct {
		base, max time.Duration
		attempt   int
	}{
		{3 * time.Second, 30 * time.Second, 1000},
		{time.Hour, math.MaxInt64, 100},
		{time.Minute, time.Second, 0}, // base above max is capped
	}
	for _, c := range cases {
		us := New("token", "", "", WithRetryBackoff(c.base, c.max))
		d := us.nextBackoff(c.attempt, 0, false)
		if d < c.max/2 || d > c.max {
			t.Errorf("nextBackoff(base=%v, max=%v, attempt=%d) = %v, want in [max/2, max]", c.base, c.max, c.attempt, d)
		}
	}
}

func TestDeleteFilesByEntityIdQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.Method != http.MethodDelete || q.Get("entityType") != "a&b" || q.Get("entityId") != "42" {
			t.Errorf("request = %s %s, want DELETE with entityType=a&b and entityId=42", r.Method, r.URL)
		}
	}))
	defer server.Close()
	if _, err := New("token", "", server.URL).DeleteFilesByEntityId(context.Background(), "a&b", 42); err != nil {
		t.Fatalf("DeleteFilesByEntityId() error = %v", err)
	}
}
