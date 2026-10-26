*** Settings ***
Documentation       Signing in, keeping a session alive, and ending it.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
The right password signs a user in
    Given a registered visitor
    When they sign in
    Then they receive a fresh token pair
    And their identity is "${EMAIL}"

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

Logging out revokes the session
    Given a signed-in platform user
    When they log out
    Then the request succeeds
    And their access token no longer works
    And their refresh token no longer works
