*** Settings ***
Documentation       A project brings its own Google/GitHub app; visitors discover
...                 and start social logins with it. The provider side is not
...                 exercised: the scenarios stop at the redirect to GitHub.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
An owner configures GitHub login
    Given a signed-in user who owns a project
    When they configure GitHub login for the project
    Then the provider is configured without exposing the secret
    And visitors are offered GitHub login

A provider is configured only once per project
    Given a signed-in user who owns a project
    And the project offers GitHub login
    When they configure GitHub login for the project
    Then the request is rejected as a conflict

A visitor is sent to GitHub with the project's app
    Given a signed-in user who owns a project
    And the project offers GitHub login
    When a visitor starts a GitHub login on the project
    Then they are sent to GitHub with client id "e2e-client-id"

A disabled provider is hidden until enabled again
    Given a signed-in user who owns a project
    And the project offers GitHub login
    When they disable GitHub login
    Then the request succeeds
    And visitors are offered no social login
    When they enable GitHub login again
    Then the request succeeds
    And visitors are offered GitHub login

Changing the client id changes where visitors are sent
    Given a signed-in user who owns a project
    And the project offers GitHub login
    When they change the GitHub client id to "e2e-rotated-id"
    And a visitor starts a GitHub login on the project
    Then they are sent to GitHub with client id "e2e-rotated-id"

Removing a provider stops offering it
    Given a signed-in user who owns a project
    And the project offers GitHub login
    When they remove GitHub login
    Then the request succeeds
    And visitors are offered no social login

The platform has no GitHub login to start
    When a visitor starts a GitHub login on the platform
    Then the request is rejected as invalid

Discovering logins for an unknown project finds nothing
    When someone asks which logins a project that does not exist offers
    Then nothing is found

Project users cannot manage login providers
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they list the project's login providers
    Then the request is forbidden
