*** Settings ***
Documentation       The admin area: organizations, their projects and members.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A user creates an organization through the app
    Given a registered visitor
    And they are signed in to the IdentityX app
    When they create the organization "UI Labs" through the app
    Then they are told "Namespace created successfully"
    And the organization's card shows up

An owner creates a project through the app
    Given a signed-in user who owns an organization
    And they are signed in to the IdentityX app
    When they open their organization in the app
    And they create the project "UI Checkout" through the app
    Then they are told "Project created successfully"
    And the project's card shows up

An owner adds a member through the app
    Given a signed-in user who owns an organization
    And another platform user
    And they are signed in to the IdentityX app
    When they add the other user to the organization through the app
    Then they are told "Member added successfully"
    And the organization shows 2 members
    And the other user is listed as a "member"

Project members show when they joined
    Given a signed-in user who owns a project
    And they are signed in to the IdentityX app
    When they open their project in the app
    And they open the project's "members" tab
    Then membership dates look real
