*** Settings ***
Documentation       Driving the IdentityX front (front/identityx) with the Browser
...                 library, which is Playwright underneath. Auth forms come from
...                 @trieoh/identityx-sdk-ts-react (react-hook-form gives every
...                 input a stable `name`); admin forms are the app's FormModal,
...                 whose inputs carry the field name too. Buttons are targeted
...                 by their aria-label, which ShadowButton sets from its text.
# Assertions retry for as long as actions wait: route guards and toasts land
# after client-side work, well past Browser's 1s default.
Library             Browser    timeout=10s    retry_assertions_for=10s
Resource            ../common.resource


*** Variables ***
${EMAIL_INPUT}              css=input[name="email"]
${PASSWORD_INPUT}           css=input[name="password"]
${CONFIRM_PASSWORD_INPUT}   css=input[name="confirmPassword"]
${SUBMIT_BUTTON}            css=form button[type="submit"]
${SIGN_UP_TOGGLE}           text=Cadastre-se
${TOAST}                    css=[data-sonner-toast]
${DIALOG}                   css=[role="dialog"]
${DIALOG_SUBMIT}            css=[role="dialog"] button[type="submit"]
${LOGOUT_BUTTON}            css=button:has-text("Logout") >> nth=0
${TEMPLATE_SUBJECT_INPUT}   css=input[placeholder="Ex.: Confirme seu e-mail"]


*** Keywords ***
# ── browser lifecycle ─────────────────────────────────────────────────────

Open The Browser
    New Browser    chromium    headless=${HEADLESS}    slowMo=${SLOW_MO}

Open A Fresh Browser Context
    [Documentation]    One context per scenario: no cookies or session leak between tests.
    New Context    locale=pt-BR    viewport={'width': 1280, 'height': 800}

# ── auth pages ────────────────────────────────────────────────────────────

Open The Sign-In Page
    New Page    ${IDX_APP_URL}/auth
    Wait For Elements State    ${EMAIL_INPUT}    visible

Open App Path
    [Documentation]    Any path of the app, in a new page of the current context.
    [Arguments]    ${path}
    New Page    ${IDX_APP_URL}${path}

Follow Link
    [Documentation]    Opens a full URL (e.g. one from an email) in a new page.
    [Arguments]    ${url}
    New Page    ${url}

Switch To The Sign-Up Form
    Click    ${SIGN_UP_TOGGLE}
    Wait For Elements State    ${CONFIRM_PASSWORD_INPUT}    visible

Fill In The Sign-In Form
    [Arguments]    ${email}    ${password}
    Fill Text    ${EMAIL_INPUT}    ${email}
    Fill Secret    ${PASSWORD_INPUT}    $password

Fill In The Sign-Up Form
    [Arguments]    ${email}    ${password}    ${confirmation}=${None}
    ${confirmation}=    Set Variable If    $confirmation is None    ${password}    ${confirmation}
    Fill Text    ${EMAIL_INPUT}    ${email}
    Fill Secret    ${PASSWORD_INPUT}    $password
    Fill Secret    ${CONFIRM_PASSWORD_INPUT}    $confirmation

Fill In The Email Form
    [Documentation]    The single-field forms: forgot password, resend verification.
    [Arguments]    ${email}
    Wait For Elements State    ${EMAIL_INPUT}    visible
    Fill Text    ${EMAIL_INPUT}    ${email}

Fill In The New Password Form
    [Arguments]    ${password}
    Wait For Elements State    ${CONFIRM_PASSWORD_INPUT}    visible
    Fill Secret    ${PASSWORD_INPUT}    $password
    Fill Secret    ${CONFIRM_PASSWORD_INPUT}    $password

Submit The Form
    Click    ${SUBMIT_BUTTON}

Log Out Of The App
    Click    ${LOGOUT_BUTTON}

# ── admin area ────────────────────────────────────────────────────────────

Press The Button
    [Documentation]    A page-level (non-submit) button, by its label.
    [Arguments]    ${label}
    Click    css=button[type="button"][aria-label="${label}"] >> nth=0
    Wait For Elements State    ${DIALOG}    visible

Fill In The Dialog Field
    [Arguments]    ${field}    ${value}
    Fill Text    ${DIALOG} input[name="${field}"]    ${value}

Submit The Dialog
    Click    ${DIALOG_SUBMIT}

The Dialog Should Show
    [Arguments]    ${text}
    Wait For Elements State    ${DIALOG} >> text=${text}    visible

Open The Card Called
    [Arguments]    ${name}
    Click    css=a >> text=${name}

Open The Project Tab
    [Documentation]    Project tabs link to /admin/projects/<id>/<tab>?organizationID=...
    [Arguments]    ${tab}
    Click    css=a[href*="/${tab}?"] >> nth=0

Replace The Template Subject
    [Arguments]    ${subject}
    Fill Text    ${TEMPLATE_SUBJECT_INPUT}    ${subject}
    Click    css=button[aria-label="Salvar"]

# ── assertions ────────────────────────────────────────────────────────────

Current Path Should Start With
    [Documentation]    Browser's assertion retries until the library timeout.
    [Arguments]    ${path}
    Get Url    *=    ${IDX_APP_URL}${path}

A Toast Should Show
    [Arguments]    ${text}
    Wait For Elements State    ${TOAST} >> text=${text}    visible

The Page Should Show
    [Arguments]    ${text}
    Wait For Elements State    text=${text} >> nth=0    visible

The Page Should Not Show
    [Arguments]    ${text}
    Get Element Count    text=${text}    ==    0

The Heading Should Read
    [Arguments]    ${text}
    Get Text    css=h1 >> nth=0    ==    ${text}
