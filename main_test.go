package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/seblw/wp-atrd-task/server"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

type TestSuite struct {
	host string
	res  *http.Response
}

func (ts *TestSuite) iSendARequestTo(method, endpoint string) error {
	var req, err = http.NewRequest(method, ts.host+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request %w", err)
	}

	ts.res = res
	return nil
}

func (ts *TestSuite) iSendARequestToWith(method, endpoint, body string) error {
	var req, err = http.NewRequest(method, ts.host+endpoint, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request %w", err)
	}

	ts.res = res
	return nil
}

func (ts *TestSuite) theJSONResponseShouldContainSecretData() error {
	body, err := ioutil.ReadAll(ts.res.Body)
	defer ts.res.Body.Close()
	if err != nil {
		return err
	}

	var secret server.Secret
	json.Unmarshal(body, &secret)

	if secret.SecretText == "" {
		return fmt.Errorf("response doesn't contain secretText")
	}

	return nil
}

func (ts *TestSuite) theResponseCodeShouldBe(code int) error {
	if ts.res.StatusCode != code {
		return fmt.Errorf("response status code doesn't match. Got: %d, expected: %d", ts.res.StatusCode, code)
	}
	return nil

}

func FeatureContext(s *godog.Suite) {
	ts := TestSuite{
		host: "http://localhost:8080",
	}

	s.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, ts.iSendARequestTo)
	s.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, ts.iSendARequestToWith)
	s.Step(`^the JSON response should contain secret data$`, ts.theJSONResponseShouldContainSecretData)
	s.Step(`^the response code should be (\d+)$`, ts.theResponseCodeShouldBe)
}
