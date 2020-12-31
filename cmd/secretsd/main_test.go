package main

import (
		"context"
		"encoding/json"
		"fmt"
		"github.com/cucumber/godog"
		"github.com/mkalafior/wp-atrd-task/internal/mongo"
		"go.mongodb.org/mongo-driver/bson"
		"io/ioutil"
		"net/http"
		"strconv"
		"strings"
)

type TestSuit struct {
		res *http.Response
}

func (ts *TestSuit) iSendARequestTo(method, uri string) error {
		var req, err = http.NewRequest(method, "http://localhost:3000"+uri, nil)
		if err != nil {
				return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
				return err
		}

		ts.res = res

		return nil
}

func (ts *TestSuit) iSendARequestToWith(method, uri, payload string) error {
		var req, err = http.NewRequest(method, "http://localhost:3000"+uri, strings.NewReader(payload))
		if err != nil {
				return err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
				return err
		}

		ts.res = res
		return nil
}

func (ts *TestSuit) theJSONResponseShouldContainSecretData() error {
		var err error
		var body interface{}

		resBody, err := ioutil.ReadAll(ts.res.Body)
		defer ts.res.Body.Close()

		if err != nil {
				return err
		}

		if err = json.Unmarshal(resBody, &body); err != nil {
				return nil
		}
		mBody := body.(map[string]interface{})
		secretText := mBody["secretText"].(string)
		if secretText == "" {
				return fmt.Errorf("there is no secret in the response %v", body)
		}

		return nil
}

func (ts *TestSuit) theResponseCodeShouldBe(code int) error {
		if code != ts.res.StatusCode {
				return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, ts.res.StatusCode)
		}
		return nil
}

func (ts *TestSuit) thereIsSecret(secrets *godog.Table) error {
		db := mongo.NewDb("mongodb://root:root@localhost:27017")
		defer mongo.Close()

		for i := 0; i < len(secrets.Rows); i++ {
				expirationdate, _ := strconv.Atoi(secrets.Rows[i].Cells[2].Value)
				createdat, _ := strconv.Atoi(secrets.Rows[i].Cells[4].Value)
				views, _ := strconv.Atoi(secrets.Rows[i].Cells[3].Value)
				secret := bson.M{
						"hash":           secrets.Rows[i].Cells[0].Value,
						"secret":         secrets.Rows[i].Cells[1].Value,
						"expirationdate": expirationdate,
						"viewsleft":      views,
						"createdat":      createdat,
				}
				_, err := db.Collection("secrets").InsertOne(context.TODO(), secret)
				if err != nil {
						return err
				}
		}

		return nil
}

func FeatureContext(s *godog.Suite) {
		ts := TestSuit{}
		s.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, ts.iSendARequestTo)
		s.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, ts.iSendARequestToWith)
		s.Step(`^the JSON response should contain secret data$`, ts.theJSONResponseShouldContainSecretData)
		s.Step(`^the response code should be (\d+)$`, ts.theResponseCodeShouldBe)
		s.Step(`^there is a secret:$`, ts.thereIsSecret)

}
