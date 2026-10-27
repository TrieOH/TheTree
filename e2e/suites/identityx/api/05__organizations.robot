*** Settings ***
Documentation       Platform users build tenants: organizations, their projects,
...                 and API keys for server-to-server access.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A platform user creates an organization they own
    Given a signed-in platform user
    When they create an organization called "E2E Labs"
    Then the organization is created
    And they own the organization
    And the organization shows in their organizations

An organization owner creates a project in it
    Given a signed-in user who owns an organization
    When they create a project called "E2E Checkout" in their organization
    Then the project belongs to their organization

An API key authenticates as the user who made it
    Given a signed-in user who owns a project
    When they create an API key for the project
    Then the raw API key is handed over
    And the API key authenticates as them

Project users cannot create organizations
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they try to create an organization
    # A project token is the wrong realm, not too few rights: IdentityX answers 401.
    Then the request is unauthorized
