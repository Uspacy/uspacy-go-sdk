package api_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Uspacy/uspacy-go-sdk/v2/api"
	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
)

// These examples compile with the tests but do not run: they need a real Uspacy portal.

func ExampleNew() {
	client := api.New("<access token>", "<refresh token>", "https://example.uspacy.ua",
		api.WithMaxRetries(5),
		api.WithRetryBackoff(time.Second, 20*time.Second),
		api.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	leads, err := client.GetLeads(ctx, url.Values{"page": {"1"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(leads.Data))
}

func ExampleWithHeader() {
	client := api.New("<access token>", "<refresh token>", "https://example.uspacy.ua")

	id, _, err := client.CreateEntity(context.Background(), crm.LeadsNum.GetUrl(),
		map[string]any{"title": "New lead"},
		api.WithHeader("X-Request-Id", "3f1c9a52"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func ExampleHTTPError() {
	client := api.New("<access token>", "<refresh token>", "https://example.uspacy.ua")

	_, err := client.GetEntity(context.Background(), crm.LeadsNum.GetUrl(), 42)

	var httpErr *api.HTTPError
	switch {
	case errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound:
		fmt.Println("lead not found")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("timed out")
	case err != nil:
		log.Fatal(err)
	}
}

func ExampleUspacy_Tokens() {
	client := api.New("<access token>", "<refresh token>", "https://example.uspacy.ua")

	// ... calls that may refresh the tokens on a 401 ...

	access, refresh := client.Tokens()
	fmt.Println(access != "", refresh != "") // store both, e.g. in a database
}
