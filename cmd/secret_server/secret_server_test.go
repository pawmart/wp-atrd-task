package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pawmart/wp-atrd-task/internal/http/app"
	"github.com/pawmart/wp-atrd-task/internal/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type apiFeature struct {
	resp *httptest.ResponseRecorder
	app  *app.App
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.resp = httptest.NewRecorder()
}

func testMain(m *testing.M) {
	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}

	godog.BindFlags("godog.", flag.CommandLine, &opts)

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opts)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func (a *apiFeature) iSendARequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return
	}

	a.makeRequest(req)

	return nil
}

func (a *apiFeature) makeRequest(r *http.Request) {
	a.app.NewRoutes().ServeHTTP(a.resp, r)
}

func (a *apiFeature) iSendARequestToWith(method, endpoint, data string) (err error) {
	req, err := http.NewRequest(method, endpoint, strings.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	a.makeRequest(req)

	return nil
}

func (a *apiFeature) theJSONResponseShouldContainSecretData() (err error) {
	var data map[string]interface{}

	if err = json.Unmarshal(a.resp.Body.Bytes(), &data); err != nil {
		return
	}

	value, exist := data["secretText"]
	if !exist || value == "" {
		err = fmt.Errorf("secretText shouldnt be empty")
		return
	}

	return nil
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

func FeatureContext(ctx *godog.Suite) {
	id, _ := uuid.Parse("b75ce598-f349-4c61-9246-2053e230187d")
	m := map[uuid.UUID]storage.Secret{
		id: {
			Id:             id,
			Value:          "test",
			CreatedAt:      time.Now(),
			ExpiresAfter:   nil,
			RemainingViews: 3,
		},
	}

	gin.SetMode(gin.TestMode)
	api := &apiFeature{resp: httptest.NewRecorder(), app: app.NewApp(storage.NewInMemoryStorageWithData(m))}

	ctx.BeforeScenario(api.resetResponse)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, api.iSendARequestTo)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, api.iSendARequestToWith)
	ctx.Step(`^the JSON response should contain secret data$`, api.theJSONResponseShouldContainSecretData)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
}
