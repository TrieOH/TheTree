*** Settings ***
Documentation       Signing in, up and out through the IdentityX app.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A registered user signs in and lands on the admin area
    Given a registered visitor
    When they sign in through the IdentityX app
    Then they land on the admin area

A wrong password keeps the user out
    Given a registered visitor
    When they sign in through the IdentityX app with a wrong password
    Then they see a warning saying "Auth Failed"

A visitor creates an account through the app
    Given a visitor with a strong password
    When they sign up through the IdentityX app
    Then a verification email arrives in their inbox

A mistyped password confirmation creates no account
    Given a visitor with a strong password
    When they sign up through the IdentityX app confirming a different password
    Then no verification email is sent to them

A guest opening the admin area is asked to sign in
    When they open "/admin" in the app
    Then they are asked to sign in

Signing in returns the user to the page they asked for
    Given a registered visitor
    When they open "/admin/introspect" in the app
    And they sign in on the page they were sent to
    Then they land on "/admin/introspect"

Logging out closes the admin area
    Given a registered visitor
    And they are signed in to the IdentityX app
    When they log out of the app
    And they open "/admin" in the app
    Then they are asked to sign in

The introspect page shows who is signed in
    Given a registered visitor
    And they are signed in to the IdentityX app
    When they open "/admin/introspect" in the app
    Then their email is shown
