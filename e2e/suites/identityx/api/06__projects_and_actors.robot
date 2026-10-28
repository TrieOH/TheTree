*** Settings ***
Documentation       Projects inside an organization: who is in them, and the
...                 actors (the project's own users and services) they hold.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
An organization owner creates a project in it
    Given a signed-in user who owns an organization
    When they create a project called "E2E Checkout" in their organization
    Then the project belongs to their organization

The new project is listed in its organization
    Given a signed-in user who owns a project
    When they list the organization's projects
    Then the project is listed

The project's creator is its owner
    Given a signed-in user who owns a project
    Then they are listed as the project's owner

Brand slugs take lowercase letters only
    Given a signed-in user who owns an organization
    When they create a project with the brand slug "Not-Valid1"
    Then the request is rejected as invalid

Brand slugs are unique
    Given a signed-in user who owns a project
    When they create another project with the same brand slug
    Then the request is rejected as a conflict

Two projects cannot share a domain
    Given a signed-in user who owns a project
    When they create another project on the same domain
    Then the request is rejected as a conflict

An owner adds and removes a project member
    Given a signed-in user who owns a project
    And another platform user
    When they add the other user to the project as "member"
    Then the member is added
    And the other user is a member of the project
    When they remove the other user from the project
    Then the request succeeds
    And the other user is not a member of the project

A project user is one of the project's actors
    Given a signed-in user who owns a project
    And a user signed in to that project
    When the owner lists the project's actors
    Then the user is among the actors

The owner looks a project user up by email
    Given a signed-in user who owns a project
    And a user signed in to that project
    When the owner looks the user up by email
    Then the actor found is the user

The owner looks a project user up by id
    Given a signed-in user who owns a project
    And a user signed in to that project
    When the owner looks the user up by id
    Then the actor found is the user

Looking up an unknown actor finds nothing
    Given a signed-in user who owns a project
    When the owner looks up an actor that does not exist
    Then nothing is found

Outsiders cannot list a project's actors
    Given a signed-in user who owns a project
    When an outsider lists the project's actors
    Then the request is forbidden

The owner provisions a service actor
    Given a signed-in user who owns a project
    When the owner provisions a service actor
    Then a service actor of the project is created

A project user's identity is scoped to the project
    Given a signed-in user who owns a project
    And a user signed in to that project
    Then their identity belongs to that project
    And they are identified as a human using a token
