*** Settings ***
Documentation       Platform users build tenants: organizations, the people in
...                 them, and the projects they hold.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A platform user creates an organization they own
    Given a signed-in platform user
    When they create an organization called "E2E Labs"
    Then the organization is created
    And they own the organization
    And the organization shows in their organizations

The creator is listed as the organization's owner
    Given a signed-in user who owns an organization
    When they list the organization's members
    Then they are listed as its owner

Organization slugs are unique
    Given a signed-in user who owns an organization
    When they create another organization with the same slug
    Then the request is rejected as a conflict

An organization needs a real name and slug
    Given a signed-in platform user
    When they create an organization with a two-letter name and a one-letter slug
    Then the request is rejected as invalid
    And the error points at the "name" field
    And the error points at the "slug" field

Project users cannot create organizations
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they try to create an organization
    # A project token is the wrong realm, not too few rights: IdentityX answers 401.
    Then the request is unauthorized

An owner adds a platform user to the organization
    Given a signed-in user who owns an organization
    And another platform user
    When they add the other user to the organization as "admin"
    Then the member is added
    And the other user is listed as a "admin"
    And the organization shows in the other user's organizations

Adding the same member twice is a conflict
    Given a signed-in user who owns an organization
    And another platform user
    And the other user is a "member" of the organization
    When they add the other user to the organization again
    Then the request is rejected as a conflict

Only registered users can be added
    Given a signed-in user who owns an organization
    When they add an email nobody registered to the organization
    Then nothing is found

Removing a member takes their access away
    Given a signed-in user who owns an organization
    And another platform user
    And the other user is a "member" of the organization
    When they remove the other user from the organization
    Then the request succeeds
    And the other user can no longer see the organization's members

Outsiders cannot see an organization's members
    Given a signed-in user who owns an organization
    And another platform user
    And the other user takes over
    When they list the organization's members
    Then the request is forbidden

Outsiders cannot see an organization's projects
    Given a signed-in user who owns an organization
    And another platform user
    And the other user takes over
    When they list the organization's projects
    Then the request is forbidden

Plain members cannot add people
    Given a signed-in user who owns an organization
    And another platform user
    And the other user is a "member" of the organization
    And the other user takes over
    When they try to add someone to the organization
    Then the request is forbidden

Plain members cannot create projects
    Given a signed-in user who owns an organization
    And another platform user
    And the other user is a "member" of the organization
    And the other user takes over
    When they try to create a project in the organization
    Then the request is forbidden

Admins can create projects
    Given a signed-in user who owns an organization
    And another platform user
    And the other user is a "admin" of the organization
    And the other user takes over
    When they create a project called "E2E Admin Project" in their organization
    Then the project belongs to their organization
