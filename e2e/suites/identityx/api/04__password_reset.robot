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
