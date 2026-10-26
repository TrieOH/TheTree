*** Settings ***
Documentation       Driving the IdentityX front (front/identityx) with the Browser
...                 library, which is Playwright underneath. Selectors target the
...                 sign-in/up forms from @trieoh/identityx-sdk-ts-react's
...                 ModernAuth: react-hook-form gives every input a stable `name`.
Library             Browser    timeout=10s
Resource            ../common.resource


*** Variables ***
${EMAIL_INPUT}              css=input[name="email"]
${PASSWORD_INPUT}           css=input[name="password"]
${CONFIRM_PASSWORD_INPUT}   css=input[name="confirmPassword"]
${SUBMIT_BUTTON}            css=form button[type="submit"]
${SIGN_UP_TOGGLE}           text=Cadastre-se
${TOAST}                    css=[data-sonner-toast]


*** Keywords ***
Open The Browser
    New Browser    chromium    headless=${HEADLESS}

Open A Fresh Browser Context
    [Documentation]    One context per scenario: no cookies or session leak between tests.
    New Context    locale=pt-BR    viewport={'width': 1280, 'height': 800}

Open The Sign-In Page
    New Page    ${IDX_APP_URL}/auth
    Wait For Elements State    ${EMAIL_INPUT}    visible

Switch To The Sign-Up Form
    Click    ${SIGN_UP_TOGGLE}
    Wait For Elements State    ${CONFIRM_PASSWORD_INPUT}    visible

Fill In The Sign-In Form
    [Arguments]    ${email}    ${password}
    Fill Text    ${EMAIL_INPUT}    ${email}
    Fill Secret    ${PASSWORD_INPUT}    $password

Fill In The Sign-Up Form
    [Arguments]    ${email}    ${password}
    Fill Text    ${EMAIL_INPUT}    ${email}
    Fill Secret    ${PASSWORD_INPUT}    $password
    Fill Secret    ${CONFIRM_PASSWORD_INPUT}    $password

Submit The Form
    Click    ${SUBMIT_BUTTON}

Current Path Should Start With
    [Documentation]    Browser's assertion retries until the library timeout.
    [Arguments]    ${path}
    Get Url    *=    ${IDX_APP_URL}${path}

A Toast Should Show
    [Arguments]    ${text}
    Wait For Elements State    ${TOAST} >> text=${text}    visible
