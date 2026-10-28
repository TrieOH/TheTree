*** Settings ***
Documentation       A project rewrites the emails IdentityX sends its users; the
...                 action link always stays IdentityX's.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
Every email starts from the default template
    Given a signed-in user who owns a project
    When they list the project's email templates
    Then every kind of email comes with the default template

An owner customizes the verification email
    Given a signed-in user who owns a project
    When they customize the "verify" email
    Then the request succeeds
    And the "verify" email is the project's own

A customized verification email is what new users get
    Given a signed-in user who owns a project
    And the project has a customized "verify" email
    And a user signed in to that project
    Then the verification email they get uses the project's template

A template must keep the action link
    Given a signed-in user who owns a project
    When they save a "verify" email without the action link
    Then the request is rejected as invalid
    And the error says "ActionURL"

A template must render
    Given a signed-in user who owns a project
    When they save a "verify" email that does not render
    Then the request is rejected as invalid

Only known kinds of email exist
    Given a signed-in user who owns a project
    When they ask for the "welcome" email template
    Then the request is rejected as invalid

Resetting a template brings the default back
    Given a signed-in user who owns a project
    And the project has a customized "reset" email
    When they reset the "reset" email to the default
    Then the request succeeds
    And the "reset" email is back to the default

Project users cannot see the templates
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they list the project's email templates
    Then the request is forbidden

Outsiders cannot see the templates
    Given a signed-in user who owns a project
    And another platform user
    And the other user takes over
    When they list the project's email templates
    Then the request is forbidden
