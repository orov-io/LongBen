@ping
Feature: As a user when I call the ping endpoint, I would like to receive a pong response

    Scenario: Valid call
      Given I have a ping call
      When I receive the response
      Then I should receive a pong response

    Scenario: Invalid call
        Given I have an invalid call
        When I receive the response
        Then Code should be a Not Found HTTP Code