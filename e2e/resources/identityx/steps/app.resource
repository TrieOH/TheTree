*** Settings ***
Documentation       Sentences about using the IdentityX app in a browser. Built from
...                 ui.resource; set-up that is not what a scenario is about still
...                 goes through the API (see the other step files).
Library             Collections
Resource            common.resource
Resource            ../ui.resource


*** Variables ***
${APP_ORGANIZATION}     ${None}
${APP_PROJECT}          ${None}


*** Keywords ***
# ── signing in and out ────────────────────────────────────────────────────

They sign in through the IdentityX app
    Open The Sign-In Page
    Fill In The Sign-In Form    ${EMAIL}    ${PASSWORD}
    Submit The Form

They sign in through the IdentityX app with a wrong password
    Open The Sign-In Page
    Fill In The Sign-In Form    ${EMAIL}    Wr0ng-${PASSWORD}
    Submit The Form

They sign in through the IdentityX app with the password "${password}"
    Open The Sign-In Page
    Fill In The Sign-In Form    ${EMAIL}    ${password}
    Submit The Form

They are signed in to the IdentityX app
    They sign in through the IdentityX app
    They land on the admin area

They sign in on the page they were sent to
    Fill In The Sign-In Form    ${EMAIL}    ${PASSWORD}
    Submit The Form

They sign up through the IdentityX app
    Open The Sign-In Page
    Switch To The Sign-Up Form
    Fill In The Sign-Up Form    ${EMAIL}    ${PASSWORD}
    Submit The Form

They sign up through the IdentityX app confirming a different password
    Open The Sign-In Page
    Switch To The Sign-Up Form
    Fill In The Sign-Up Form    ${EMAIL}    ${PASSWORD}    Different-${PASSWORD}
    Submit The Form

They open "${path}" in the app
    Open App Path    ${path}

They log out of the app
    Log Out Of The App

They land on the admin area
    Current Path Should Start With    /admin

They land on "${path}"
    Current Path Should Start With    ${path}

They are asked to sign in
    Current Path Should Start With    /auth
    The Page Should Show    Bem-vindo de volta

They are back on the sign-in page
    They are asked to sign in

They see a warning saying "${text}"
    A Toast Should Show    ${text}

They are told "${text}"
    A Toast Should Show    ${text}

The page reads "${text}"
    The Page Should Show    ${text}

The page heading reads "${text}"
    The Heading Should Read    ${text}

The page does not mention "${text}"
    The Page Should Not Show    ${text}

# ── recovering an account ─────────────────────────────────────────────────

They ask for a password reset through the app
    Open App Path    /auth/forgot-password
    Fill In The Email Form    ${EMAIL}
    Submit The Form

They ask for a new verification email through the app
    Open App Path    /auth/resend-verification
    Fill In The Email Form    ${EMAIL}
    Submit The Form

They follow the reset link from the inbox and choose "${password}"
    ${link}=    Action Link From Email    ${EMAIL}    ${RESET_SUBJECT}
    Follow Link    ${link}
    Fill In The New Password Form    ${password}
    Submit The Form

They follow the verification link from the inbox
    ${link}=    Action Link From Email    ${EMAIL}    ${VERIFY_SUBJECT}
    Follow Link    ${link}

# ── the admin area ────────────────────────────────────────────────────────

They create the organization "${name}" through the app
    ${slug}=    Unique Slug    ui
    VAR    ${APP_ORGANIZATION}    ${name}    scope=TEST
    Press The Button    Create Organization
    Fill In The Dialog Field    name    ${name}
    Fill In The Dialog Field    slug    ${slug}
    Submit The Dialog

They open their organization in the app
    Open The Card Called    ${ORGANIZATION}[name]
    Current Path Should Start With    /admin/${ORGANIZATION}[id]

They create the project "${name}" through the app
    ${slug}=    Unique Brand Slug    ui
    ${domain}=    Unique Project Domain
    VAR    ${APP_PROJECT}    ${name}    scope=TEST
    Press The Button    Add Project
    Fill In The Dialog Field    name    ${name}
    Fill In The Dialog Field    brand_slug    ${slug}
    Fill In The Dialog Field    domain    ${domain}
    Submit The Dialog

They add the other user to the organization through the app
    Open App Path    /admin/${ORGANIZATION}[id]/members
    Press The Button    Add Member
    Fill In The Dialog Field    actor_email    ${OTHER_EMAIL}
    Submit The Dialog

They open their project in the app
    Open App Path    /admin/projects/${PROJECT}[id]?organizationID=${ORGANIZATION}[id]
    The Page Should Show    Api Keys

They open the project's "${tab}" tab
    Open The Project Tab    ${tab}

They create the API key "${name}" through the app
    Press The Button    Create API Key
    Fill In The Dialog Field    name    ${name}
    Fill In The Dialog Field    env    sandbox
    Submit The Dialog

They add the capability "${resource}:${action}" through the app
    Press The Button    Add Capability
    Fill In The Dialog Field    resource    ${resource}
    Fill In The Dialog Field    action    ${action}
    Submit The Dialog

They change the verification email subject to "${subject}" through the app
    Replace The Template Subject    ${subject}
    VAR    ${APP_SUBJECT}    ${subject}    scope=TEST

The new key is shown once
    The Dialog Should Show    API Key Created
    The Dialog Should Show    You won't be able to see it again

The organization's card shows up
    The Page Should Show    ${APP_ORGANIZATION}

The project's card shows up
    The Page Should Show    ${APP_PROJECT}

The organization shows ${count} members
    The Page Should Show    ${count} members in this organization

Their email is shown
    The Page Should Show    ${EMAIL}

Membership dates look real
    [Documentation]    joined_at must be a real date, not Go's zero time (year 1).
    The Page Should Show    Joined
    The Page Should Not Show    years ago

The verification email subject is saved for the project
    ${response}=    IDX Get Email Template    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    verify
    Response Status Should Be    200    ${response}
    ${template}=    Data Of    ${response}
    Should Be Equal    ${template}[source]    override
    Should Be Equal    ${template}[subject]    ${APP_SUBJECT}
