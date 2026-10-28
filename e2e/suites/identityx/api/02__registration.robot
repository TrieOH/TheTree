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

Opening the verification link twice is harmless
    Given a registered visitor
    When they confirm their email with the link from the inbox
    And they open the same verification link again
    Then the request succeeds

A forged verification token is rejected
    When someone verifies an email with a forged token
    Then the request is rejected as invalid

A visitor who lost the email can ask for another one
    Given a registered visitor
    When they ask for a new verification email
    Then the request succeeds
    And a second verification email arrives in their inbox

A verified account is not sent another verification email
    Given a verified user
    When they ask for a new verification email
    Then the request succeeds
    And no second verification email is sent to them

Asking to verify an unknown email gives nothing away
    When someone asks for a verification email for an address nobody registered
    Then the request succeeds
    And no verification email is sent to them

The same email cannot register twice
    Given a registered visitor
    When they register on the platform again
    Then the request is rejected as a conflict

Emails are matched regardless of case
    Given a registered visitor
    When they sign in with their email in capitals
    Then they receive a fresh token pair

A password without a symbol, number and capital is rejected
    Given a visitor with the password "weakpassword"
    When they register on the platform
    Then the request is rejected as invalid
    And the error points at the "password" field

A password longer than 72 characters is rejected
    Given a visitor with the password "Aa1!aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    When they register on the platform
    Then the request is rejected as invalid
    And the error points at the "password" field

An address that is not an email is rejected
    Given a visitor with the email "not-an-email"
    When they register on the platform
    Then the request is rejected as invalid
    And the error points at the "email" field
