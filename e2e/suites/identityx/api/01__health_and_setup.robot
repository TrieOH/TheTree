*** Settings ***
Documentation       The service is up and the platform has its super admin.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
The service reports it is healthy
    When they check the service health
    Then the service reports it is healthy

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
