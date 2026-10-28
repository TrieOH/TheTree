*** Settings ***
Documentation       Signing in, keeping a session alive, and ending it.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
The right password signs a user in
    Given a registered visitor
    When they sign in
    Then they receive a fresh token pair
    And their identity is "${EMAIL}"
    And their identity belongs to no project
    And they are identified as a human using a token
    And the access token expires before the refresh token

A wrong password is refused
    Given a registered visitor
    When they sign in with a wrong password
    Then the request is unauthorized

An unknown email is refused
    When someone signs in with an email nobody registered
    Then the request is unauthorized

Refreshing a session rotates both tokens
    Given a signed-in platform user
    When they refresh their session
    Then they receive a fresh token pair
    And the new tokens replace the old ones
    And their old refresh token no longer works
    And their old access token no longer works

An access token cannot be used to refresh a session
    Given a signed-in platform user
    When they try to refresh with their access token
    Then the request is unauthorized

A malformed refresh token is refused
    When someone refreshes with a malformed token
    Then the request is unauthorized

Logging out revokes the session
    Given a signed-in platform user
    When they log out
    Then the request succeeds
    And their access token no longer works
    And their refresh token no longer works

Introspection needs a credential
    When someone asks who they are without credentials
    Then the request is unauthorized

A tampered access token is refused
    Given a signed-in platform user
    When someone asks who they are with a forged token
    Then the request is unauthorized

A project user signs in to their project only
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they sign in to the platform instead of their project
    Then the request is unauthorized

A platform account does not exist inside a project
    Given a signed-in user who owns a project
    When they sign in to that project with their platform account
    Then the request is unauthorized
