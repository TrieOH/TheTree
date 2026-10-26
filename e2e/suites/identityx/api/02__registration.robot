*** Settings ***
Documentation       A visitor creates a platform account and confirms their email.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A visitor registers with a strong password
    Given a visitor with a strong password
    When they register on the platform
    Then the registration is accepted
    And a verification email arrives in their inbox

A new account starts unverified
    Given a registered visitor
    When they sign in
    Then they receive a fresh token pair
    And their identity shows the email as unverified

The link from the inbox verifies the email
    Given a registered visitor
    When they confirm their email with the link from the inbox
    And they sign in
    Then they receive a fresh token pair
    And their identity shows the email as verified

The same email cannot register twice
    Given a registered visitor
    When they register on the platform again
    Then the request is rejected as a conflict

A password without a symbol, number and capital is rejected
    Given a visitor with the password "weakpassword"
    When they register on the platform
    Then the request is rejected as invalid
    And the error points at the "password" field
