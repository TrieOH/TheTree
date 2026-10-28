*** Settings ***
Documentation       A project's terms of service: introduced, versioned, announced
...                 to its users, and agreed to before anyone signs up (LGPD art. 8-9).
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A project starts without terms
    Given a signed-in user who owns a project
    When someone reads the project's terms
    Then nothing is found

An owner introduces terms as version 1
    Given a signed-in user who owns a project
    When they introduce the terms "Be nice to each other."
    Then the terms are at version 1
    And the terms take effect 30 days after they are announced
    And the project's terms read "Be nice to each other."

Terms are introduced only once
    Given a signed-in user who owns a project
    And the project has the terms "Be nice to each other."
    When they introduce the terms "Be nicer."
    Then the request is rejected as a conflict

Updating the terms bumps the version
    Given a signed-in user who owns a project
    And the project has the terms "Be nice to each other."
    When they update the terms to "Be nicer to each other."
    Then the terms are at version 2
    And the project's terms read "Be nicer to each other."

Terms that do not exist cannot be updated
    Given a signed-in user who owns a project
    When they update the terms to "Be nicer."
    Then nothing is found

Existing users are told about new terms
    Given a signed-in user who owns a project
    And a user signed in to that project
    And the project has the terms "Be nice to each other."
    Then they are told about the new terms by email

Signing up requires accepting the terms
    Given a signed-in user who owns a project
    And the project has the terms "Be nice to each other."
    And a visitor with a strong password
    When they sign up to the project without accepting the terms
    Then the request is rejected as invalid

Declining the terms is not accepting them
    Given a signed-in user who owns a project
    And the project has the terms "Be nice to each other."
    And a visitor with a strong password
    When they sign up to the project declining the terms
    Then the request is rejected as invalid

Accepting the terms lets a visitor sign up
    Given a signed-in user who owns a project
    And the project has the terms "Be nice to each other."
    And a visitor with a strong password
    When they sign up to the project accepting the terms
    Then the registration is accepted
    And their acceptance of version 1 is on record

An existing user accepts updated terms
    Given a signed-in user who owns a project
    And a user signed in to that project
    And the project has the terms "Be nice to each other."
    When they accept the current terms
    Then the request succeeds
    And their acceptance of version 1 is on record

Accepting the same terms twice records one acceptance
    Given a signed-in user who owns a project
    And a user signed in to that project
    And the project has the terms "Be nice to each other."
    When they accept the current terms
    And they accept the current terms again
    Then the request succeeds
    And they have accepted the terms only once

Project users cannot read the consent ledger
    Given a signed-in user who owns a project
    And a user signed in to that project
    And the project has the terms "Be nice to each other."
    When they read the terms acceptances
    Then the request is forbidden
