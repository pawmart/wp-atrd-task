Feature: Fetch a secret
  As an API user
  I want to get a secret
  So I can get details of the secret

  Scenario: Get secret
    When I send a "GET" request to "/v1/secret/b75ce598-f349-4c61-9246-2053e230187d"
    And the JSON response should contain secret data
    Then the response code should be 200

  Scenario: Fail to get a secret
    When I send a "GET" request to "/v1/secret/f3c5f34a-3985-44b2-bb1d-a51ffda32baf"
    Then the response code should be 404