package main

import "github.com/cucumber/godog"

func iSendARequestTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func iSendARequestToWith(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func theJSONResponseShouldContainSecretData() error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, iSendARequestTo)
	s.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, iSendARequestToWith)
	s.Step(`^the JSON response should contain secret data$`, theJSONResponseShouldContainSecretData)
	s.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
}
