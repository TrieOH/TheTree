*** Settings ***
Documentation       Project profiles: a JSON Schema the owner defines, and the
...                 profiles users fill in against it.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A project starts without a profile schema
    Given a signed-in user who owns a project
    When someone reads the project's profile schema
    Then nothing is found

An owner publishes a profile schema
    Given a signed-in user who owns a project
    When they publish a profile schema requiring a full name
    Then the profile schema is at version 1
    And anyone can read the schema

A user fills in their own profile
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    When they set their full name to "Jane Doe" with a handle
    Then the request succeeds
    And their public profile shows the full name "Jane Doe"

A profile is found by its handle
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And they have the full name "Jane Doe" and a handle
    When someone looks the profile up by its handle
    Then the profile found is theirs

An unknown handle finds nothing
    When someone looks up a handle nobody has
    Then nothing is found

Editing a profile without a handle keeps the handle
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And they have the full name "Jane Doe" and a handle
    When they change their full name to "Jane Roe"
    Then the request succeeds
    And their public profile shows the full name "Jane Roe"
    And their profile keeps its handle

A profile missing a required field is rejected
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    When they save a profile without a full name
    Then the request is rejected as invalid

A profile with the wrong type is rejected
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    When they save a profile whose full name is a number
    Then the request is rejected as invalid

Users cannot edit each other's profiles
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And another user of that project
    When they try to edit the other user's profile
    Then the request is forbidden

Handles are unique
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And they have the full name "Jane Doe" and a handle
    And another user of that project
    And the other user takes over
    When they try to take that handle
    Then the request is rejected as a conflict

The owner can edit any profile
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And they have the full name "Jane Doe" and a handle
    When the owner renames the user to "Jane Admin"
    Then the request succeeds
    And their public profile shows the full name "Jane Admin"

A stricter schema flags existing profiles as outdated
    Given a signed-in user who owns a project
    And the project requires a full name on profiles
    And a user signed in to that project
    And they have the full name "Jane Doe" and a handle
    When the owner starts requiring an age as well
    And the owner lists the outdated profiles
    Then their profile is flagged as outdated
    And their profile is in the list

Project users cannot list outdated profiles
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they list the outdated profiles
    Then the request is forbidden
