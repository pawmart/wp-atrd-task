Feature: Create a secret
  As an API user
  I want to store a secret
  So I can later fetch it from a secure place

  Scenario: Create
    When I send a "POST" request to "/v1/secret" with "secret=asdfasdfasdfasdf&expireAfterViews=5&expireAfter=60"
    Then the response code should be 200
    And the JSON response should contain secret data
