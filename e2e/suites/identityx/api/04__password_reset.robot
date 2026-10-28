*** Settings ***
Documentation       A user who forgot their password gets back in by email.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A forgotten password is reset through the inbox
    Given a registered visitor
    When they ask for a password reset
    And they choose the new password "N3w-Passw0rd!" with the link from the inbox
    Then the request succeeds
    And they can sign in with the password "N3w-Passw0rd!"
    And they can no longer sign in with the password "${STRONG_PASSWORD}"

A reset link works only once
    Given a registered visitor
    When they ask for a password reset
    And they choose the new password "N3w-Passw0rd!" with the link from the inbox
    And the same reset link is used again for "An0ther-Passw0rd!"
    Then the request is rejected as invalid
    And they can sign in with the password "N3w-Passw0rd!"

The new password must be strong too
    Given a registered visitor
    When they ask for a password reset
    And they choose the new password "weakpassword" with the link from the inbox
    Then the request is rejected as invalid
    And the error points at the "password" field
    And they can sign in with the password "${STRONG_PASSWORD}"

A verification link cannot reset a password
    Given a registered visitor
    When they use the verification link to reset their password
    Then the request is rejected as invalid
    And they can sign in with the password "${STRONG_PASSWORD}"

A forged reset token is rejected
    When someone resets a password with a forged token
    Then the request is rejected as invalid

Asking to reset an unknown email gives nothing away
    When someone asks for a password reset for an address nobody registered
    Then the request succeeds
    And no password reset email is sent to them

Resetting the password ends existing sessions
    Given a signed-in platform user
    When they ask for a password reset
    And they choose the new password "N3w-Passw0rd!" with the link from the inbox
    Then the request succeeds
    And their refresh token no longer works
