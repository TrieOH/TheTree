*** Settings ***
Documentation       The service is up, describes itself, publishes its signing keys,
...                 and the platform has its super admin.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
The service reports it is healthy
    When they check the service health
    Then the service reports it is healthy

The service serves its own API specification
    When someone fetches the API specification
    Then the specification describes the IdentityX API

Setup is complete once the first admin exists
    When someone asks whether setup is still pending
    Then setup is reported as complete

The platform cannot be set up twice
    When someone tries to set up the platform again
    Then the request is rejected as a conflict

The platform admin can sign in
    Given the platform admin
    When they sign in
    Then they receive a fresh token pair
    And their identity is "${ADMIN_EMAIL}"

The platform publishes the keys its tokens are signed with
    Given a signed-in platform user
    When someone fetches the platform signing keys
    Then they get Ed25519 public keys
    And their access token was signed with one of those keys

Each project signs with its own keys
    Given a signed-in user who owns a project
    And a user signed in to that project
    When someone fetches that project's signing keys
    Then they get Ed25519 public keys
    And the keys differ from the platform's
    And their access token was signed with one of those keys

The platform offers no social login until one is configured
    When someone asks which logins the platform offers
    Then no social login is offered
