package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"notsosecretsercet/pkg/adding"
	"notsosecretsercet/pkg/listing"
	"strings"
	"testing"
)

func setup(e error) *httptest.Server {
	repo := mocStore{storeErr: e}

	as := adding.NewService(&repo)
	ls := listing.NewService(&repo)
	return httptest.NewServer(Handler(as, ls))
}

func teardown(s *httptest.Server) {
	s.Close()
}

func TestAddSecret_MethodNotAllowed(t *testing.T) {
	srv := setup(nil)

	urlPostfix := apiVersion + "/secret"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}

	client := srv.Client()

	for _, method := range methods {

		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected to get 405 on method: %s, got %d", method, resp.StatusCode)
		}

	}

	teardown(srv)
}

func TestAddSecret_InvalidInput(t *testing.T) {
	srv := setup(nil)

	wrongSecretdata := []url.Values{
		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"10"},
			"expireAfter":      {"-10"},
		},

		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"0"},
			"expireAfter":      {"10"},
		},

		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"-100"},
			"expireAfter":      {"10"},
		},

		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"broke end point pls"},
			"expireAfter":      {"10"},
		},

		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"7"},
			"expireAfter":      {"broke endpoint"},
		},
		{
			"secret":           {"Test incorrect secret"},
			"expireAfterViews": {"7"},
			"expireAfter":      {"7"},
			"radnomField":      {"bambo"},
		},
	}

	urlPostfix := apiVersion + "/secret"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	for _, data := range wrongSecretdata {
		resp, err := http.PostForm(url, data)

		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", resp.StatusCode)
		}

		bytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Fatal(err)
		}

		respBody := strings.Trim(string(bytes), "\n")
		if respBody != InvalidInputMessage {
			t.Fatalf("expected  response body to be %s, got %s", InvalidInputMessage, string(bytes))
		}
	}

	teardown(srv)
}
func TestAddSecret_InternalServerError(t *testing.T) {
	srv := setup(errors.New("Test error"))

	data := url.Values{
		"secret":           {"Test correct secret"},
		"expireAfterViews": {"10"},
		"expireAfter":      {"10"},
	}

	urlPostfix := apiVersion + "/secret"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)
	resp, err := http.PostForm(url, data)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}

	teardown(srv)
}
func TestAddSecret_ValidInput(t *testing.T) {
	srv := setup(nil)

	data := url.Values{
		"secret":           {"Test correct secret"},
		"expireAfterViews": {"10"},
		"expireAfter":      {"10"},
	}

	urlPostfix := apiVersion + "/secret"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)
	res, err := http.PostForm(url, data)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	teardown(srv)
}

func TestGetSecret_MethodNotAllowed(t *testing.T) {
	srv := setup(nil)

	urlPostfix := apiVersion + "/secret/randomHash"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}

	client := srv.Client()

	for _, method := range methods {

		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected to get 405 on method: %s, got %d", method, res.StatusCode)
		}

	}

	teardown(srv)
}
func TestGetSecret_NotFound(t *testing.T) {
	srv := setup(listing.ErrNotFound)

	urlPostfix := apiVersion + "/secret/randomHash"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	res, err := http.Get(url)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
	teardown(srv)
}

func TestGetSecret_InternalServerError(t *testing.T) {
	srv := setup(errors.New("Test error"))

	urlPostfix := apiVersion + "/secret/randomHash"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	res, err := http.Get(url)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", res.StatusCode)
	}

	teardown(srv)
}
func TestGetSecret_ProperContentType(t *testing.T) {
	srv := setup(nil)

	urlPostfix := apiVersion + "/secret/randomHash"
	url := fmt.Sprintf("%s%s", srv.URL, urlPostfix)

	supportedContentTypes := []string{
		ContentTypeApplicationJSON,
		ContentTypeApplicationXML,
	}

	client := srv.Client()

	for _, contentType := range supportedContentTypes {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Accept", contentType)

		if err != nil {
			t.Fatal(err)
		}

		res, _ := client.Do(req)

		if resContentType := res.Header.Get("Content-Type"); resContentType != contentType {
			t.Fatalf("Expected content type to be %s got %s", contentType, resContentType)
		}
	}

	teardown(srv)
}

// moced db
var _ adding.Repository = &mocStore{}
var _ listing.Repository = &mocStore{}

type mocStore struct {
	storeErr error
}

func (ms *mocStore) AddSecret(s adding.Secret) (*listing.Secret, error) {
	return nil, ms.storeErr
}

func (ms *mocStore) GetSecret(hash string) (*listing.Secret, error) {
	return nil, ms.storeErr
}
