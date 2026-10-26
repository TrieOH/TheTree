*** Settings ***
Documentation       Signing in and up through the IdentityX app.
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
