*** Settings ***
Documentation       Getting back in through the app: verification and reset links
...                 from the inbox, opened in the browser.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
A forgotten password is reset through the app
    Given a registered visitor
    When they ask for a password reset through the app
    Then they are told "Link de redefinição enviado"
    When they follow the reset link from the inbox and choose "N3w-Passw0rd!"
    Then they are back on the sign-in page
    When they sign in through the IdentityX app with the password "N3w-Passw0rd!"
    Then they land on the admin area

The verification link confirms the email in the app
    Given a registered visitor
    When they follow the verification link from the inbox
    Then the page heading reads "E-mail Verificado!"
    When they sign in
    Then they receive a fresh token pair
    And their identity shows the email as verified

A forged verification link fails visibly
    When they open "/auth/verify-email?token=forged" in the app
    Then the page heading reads "Falha na verificação"

A reset page without a token says the link is broken
    When they open "/auth/reset-password" in the app
    Then the page reads "Link de redefinição inválido ou incompleto."

A visitor asks the app for a new verification email
    Given a registered visitor
    When they ask for a new verification email through the app
    Then they are told "E-mail de verificação enviado!"
    And a second verification email arrives in their inbox
