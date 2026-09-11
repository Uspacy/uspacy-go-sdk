package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
