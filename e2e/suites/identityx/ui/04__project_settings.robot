*** Settings ***
Documentation       Configuring a project from its admin pages.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
An owner creates an API key through the app
    Given a signed-in user who owns a project
    And they are signed in to the IdentityX app
    When they open their project in the app
    And they create the API key "ui-key" through the app
    Then they are told "API key created successfully"
    And the new key is shown once

An owner adds a capability through the app
    Given a signed-in user who owns a project
    And they are signed in to the IdentityX app
    When they open their project in the app
    And they open the project's "capabilities" tab
    And they add the capability "orders:refund" through the app
    Then they are told "Capability created successfully"
    And the page reads "orders"

An owner rewrites the verification email subject through the app
    Given a signed-in user who owns a project
    And they are signed in to the IdentityX app
    When they open their project in the app
    And they open the project's "emails" tab
    And they change the verification email subject to "Confirme seu e-mail na UI" through the app
    Then they are told "Template salvo"
    And the verification email subject is saved for the project
