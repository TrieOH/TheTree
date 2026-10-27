*** Settings ***
Documentation       The English the IdentityX scenarios are written in. Each step
...                 is a Given/When/Then sentence (Robot drops the prefix when it
...                 matches) built from the code-level keywords in api.resource,
...                 ui.resource and mailpit.resource.
...
...                 Steps share state through test-scoped variables: the
...                 scenario's actor (${EMAIL}/${PASSWORD}), their ${TOKENS}, and
...                 the ${RESPONSE} the last "When" got back for the "Then"s.
Library             Collections
Resource            ../common.resource
Resource            ../mailpit.resource
Resource            api.resource
Resource            ui.resource


*** Variables ***
${EMAIL}                ${None}
${PASSWORD}             ${None}
${RESPONSE}             ${None}
${TOKENS}               ${None}
${OLD_TOKENS}           ${None}
${ORGANIZATION}         ${None}
${PROJECT}              ${None}
${API_KEY}              ${None}
${USED_RESET_TOKEN}     ${None}

${VERIFY_SUBJECT}       Verify your email address
${RESET_SUBJECT}        Reset your password


*** Keywords ***
# ── the stack ─────────────────────────────────────────────────────────────

The IdentityX stack is ready
    [Documentation]    Suite setup: sessions open, the API healthy, the platform set up.
    Open IdentityX Session
    Open Mailpit Session
    Wait Until Service Is Healthy    ${IDX_URL}
    The platform has been set up

The platform has been set up
    [Documentation]    `/auth/setup` succeeds once per stack: set it up if pending, else move on.
    ${status}=    IDX Get Setup Status
    IF    ${status.status_code} == 204
        ${response}=    IDX Set Up Platform    ${ADMIN_EMAIL}    ${ADMIN_PASSWORD}
        Response Status Should Be    201    ${response}
    ELSE
        Response Status Should Be    409    ${status}
    END

# ── who the scenario is about ─────────────────────────────────────────────

A visitor with a strong password
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    VAR    ${PASSWORD}    ${STRONG_PASSWORD}    scope=TEST

A visitor with the password "${password}"
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    VAR    ${PASSWORD}    ${password}    scope=TEST

A registered visitor
    A visitor with a strong password
    They register on the platform
    The registration is accepted

A verified user
    A registered visitor
    They confirm their email with the link from the inbox
    The request succeeds

A signed-in platform user
    A registered visitor
    They sign in
    They receive a fresh token pair

The platform admin
    VAR    ${EMAIL}    ${ADMIN_EMAIL}    scope=TEST
    VAR    ${PASSWORD}    ${ADMIN_PASSWORD}    scope=TEST

A signed-in user who owns an organization
    A signed-in platform user
    They create an organization called "E2E Labs"
    The organization is created

A signed-in user who owns a project
    A signed-in user who owns an organization
    They create a project called "E2E Checkout" in their organization
    The project belongs to their organization

A user signed in to that project
    [Documentation]    Swaps the scenario's actor for a project user of ${PROJECT}.
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    VAR    ${PASSWORD}    ${STRONG_PASSWORD}    scope=TEST
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    Response Status Should Be    201    ${response}
    ${response}=    IDX Login    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    VAR    ${RESPONSE}    ${response}    scope=TEST
    They receive a fresh token pair

# ── what they do (API) ────────────────────────────────────────────────────

They check the service health
    ${response}=    IDX Get Health
    VAR    ${RESPONSE}    ${response}    scope=TEST

Someone asks whether setup is still pending
    ${response}=    IDX Get Setup Status
    VAR    ${RESPONSE}    ${response}    scope=TEST

Someone tries to set up the platform again
    ${email}=    Unique Email
    ${response}=    IDX Set Up Platform    ${email}    ${STRONG_PASSWORD}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They register on the platform
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They register on the platform again
    They register on the platform

They sign in
    ${response}=    IDX Login    ${EMAIL}    ${PASSWORD}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They sign in with a wrong password
    ${response}=    IDX Login    ${EMAIL}    Wr0ng-${PASSWORD}
    VAR    ${RESPONSE}    ${response}    scope=TEST

Someone signs in with an email nobody registered
    ${email}=    Unique Email
    ${response}=    IDX Login    ${email}    ${STRONG_PASSWORD}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They refresh their session
    VAR    ${OLD_TOKENS}    ${TOKENS}    scope=TEST
    ${response}=    IDX Refresh    ${TOKENS}[refresh_token]
    VAR    ${RESPONSE}    ${response}    scope=TEST

They log out
    ${response}=    IDX Logout    ${TOKENS}[access_token]    ${TOKENS}[refresh_token]
    VAR    ${RESPONSE}    ${response}    scope=TEST

They confirm their email with the link from the inbox
    ${token}=    Action Token From Email    ${EMAIL}    ${VERIFY_SUBJECT}
    ${response}=    IDX Verify Email    ${token}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They ask for a password reset
    ${response}=    IDX Forgot Password    ${EMAIL}
    Response Status Should Be    200    ${response}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They choose the new password "${new_password}" with the link from the inbox
    ${token}=    Action Token From Email    ${EMAIL}    ${RESET_SUBJECT}
    ${response}=    IDX Reset Password    ${token}    ${new_password}
    VAR    ${RESPONSE}    ${response}    scope=TEST
    VAR    ${USED_RESET_TOKEN}    ${token}    scope=TEST

The same reset link is used again for "${new_password}"
    ${response}=    IDX Reset Password    ${USED_RESET_TOKEN}    ${new_password}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They create an organization called "${name}"
    ${slug}=    Unique Slug    org
    ${response}=    IDX Create Organization    ${TOKENS}[access_token]    ${name}    ${slug}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They try to create an organization
    They create an organization called "Not Allowed"

They create a project called "${name}" in their organization
    ${slug}=    Unique Brand Slug    prj
    # A project needs a domain: its users' verify/reset links point there.
    # Domains are unique across projects, so each gets its own *.localhost
    # subdomain, which still resolves to the front.
    ${domain}=    Evaluate    $IDX_APP_URL.replace("://", f"://{$slug}.", 1)
    ${response}=    IDX Create Organization Project
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${name}    ${slug}    domain=${domain}
    VAR    ${RESPONSE}    ${response}    scope=TEST

They create an API key for the project
    ${response}=    IDX Create API Key    ${TOKENS}[access_token]    ${PROJECT}[id]    e2e
    VAR    ${RESPONSE}    ${response}    scope=TEST

# ── what they do (browser) ────────────────────────────────────────────────

They sign in through the IdentityX app
    Open The Sign-In Page
    Fill In The Sign-In Form    ${EMAIL}    ${PASSWORD}
    Submit The Form

They sign in through the IdentityX app with a wrong password
    Open The Sign-In Page
    Fill In The Sign-In Form    ${EMAIL}    Wr0ng-${PASSWORD}
    Submit The Form

They sign up through the IdentityX app
    Open The Sign-In Page
    Switch To The Sign-Up Form
    Fill In The Sign-Up Form    ${EMAIL}    ${PASSWORD}
    Submit The Form

# ── what should happen ────────────────────────────────────────────────────

The request succeeds
    Response Should Be Successful    ${RESPONSE}

The registration is accepted
    Response Status Should Be    201    ${RESPONSE}

The request is rejected as invalid
    Response Status Should Be    400    ${RESPONSE}

The request is rejected as a conflict
    Response Status Should Be    409    ${RESPONSE}

The request is unauthorized
    Response Status Should Be    401    ${RESPONSE}

The request is forbidden
    Response Status Should Be    403    ${RESPONSE}

The error points at the "${field}" field
    ${error}=    Error Of    ${RESPONSE}
    Dictionary Should Contain Key    ${error}    fields    msg=No field errors: ${RESPONSE.text}
    ${names}=    Evaluate    [f["field"] for f in $error["fields"]]
    List Should Contain Value    ${names}    ${field}

The service reports it is healthy
    Response Status Should Be    200    ${RESPONSE}
    ${body}=    Data Of    ${RESPONSE}
    Should Be Equal    ${body}[status]    ok
    Should Be Equal    ${body}[service]    IdentityX

Setup is reported as complete
    Response Status Should Be    409    ${RESPONSE}

A verification email arrives in their inbox
    Wait For Email To    ${EMAIL}    ${VERIFY_SUBJECT}

They receive a fresh token pair
    Response Status Should Be    200    ${RESPONSE}
    ${tokens}=    Data Of    ${RESPONSE}
    Should Not Be Empty    ${tokens}[access_token]
    Should Not Be Empty    ${tokens}[refresh_token]
    VAR    ${TOKENS}    ${tokens}    scope=TEST

The new tokens replace the old ones
    Should Not Be Equal    ${TOKENS}[access_token]    ${OLD_TOKENS}[access_token]
    Should Not Be Equal    ${TOKENS}[refresh_token]    ${OLD_TOKENS}[refresh_token]

Their old refresh token no longer works
    ${response}=    IDX Refresh    ${OLD_TOKENS}[refresh_token]
    Response Status Should Be    401    ${response}

Their old access token no longer works
    ${response}=    IDX Introspect With Token    ${OLD_TOKENS}[access_token]
    Response Status Should Be    401    ${response}

Their access token no longer works
    ${response}=    IDX Introspect With Token    ${TOKENS}[access_token]
    Response Status Should Be    401    ${response}

Their refresh token no longer works
    ${response}=    IDX Refresh    ${TOKENS}[refresh_token]
    Response Status Should Be    401    ${response}

Their identity is "${email}"
    ${subject}=    Subject Behind Token    ${TOKENS}[access_token]
    Should Be Equal    ${subject}[email]    ${email}

Their identity shows the email as verified
    ${verified_at}=    Verified At    ${TOKENS}[access_token]
    Should Not Be Equal    ${verified_at}    ${None}    msg=Email is still unverified

Their identity shows the email as unverified
    ${verified_at}=    Verified At    ${TOKENS}[access_token]
    Should Be Equal    ${verified_at}    ${None}    msg=Email is already verified

They can sign in with the password "${password}"
    ${response}=    IDX Login    ${EMAIL}    ${password}
    Response Status Should Be    200    ${response}

They can no longer sign in with the password "${password}"
    ${response}=    IDX Login    ${EMAIL}    ${password}
    Response Status Should Be    401    ${response}

The organization is created
    Response Status Should Be    201    ${RESPONSE}
    ${organization}=    Data Of    ${RESPONSE}
    VAR    ${ORGANIZATION}    ${organization}    scope=TEST

They own the organization
    ${subject}=    Subject Behind Token    ${TOKENS}[access_token]
    Should Be Equal    ${ORGANIZATION}[owner_id]    ${subject}[id]

The organization shows in their organizations
    ${response}=    IDX List Organizations    ${TOKENS}[access_token]
    Response Status Should Be    200    ${response}
    ${organizations}=    Data Of    ${response}
    ${ids}=    Evaluate    [o["id"] for o in $organizations]
    List Should Contain Value    ${ids}    ${ORGANIZATION}[id]

The project belongs to their organization
    Response Status Should Be    201    ${RESPONSE}
    ${project}=    Data Of    ${RESPONSE}
    Should Be Equal    ${project}[organization_id]    ${ORGANIZATION}[id]
    VAR    ${PROJECT}    ${project}    scope=TEST

The raw API key is handed over
    Response Status Should Be    201    ${RESPONSE}
    ${created}=    Data Of    ${RESPONSE}
    Should Not Be Empty    ${created}[raw_key]
    VAR    ${API_KEY}    ${created}[raw_key]    scope=TEST

The API key authenticates as them
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Introspect With API Key    ${API_KEY}
    Response Status Should Be    200    ${response}
    ${identity}=    Data Of    ${response}
    Should Be Equal    ${identity}[subject][id]    ${me}[id]
    Should Be Equal    ${identity}[cred][type]    api_key

They land on the admin area
    Current Path Should Start With    /admin

They see a warning saying "${text}"
    A Toast Should Show    ${text}

# ── helpers ───────────────────────────────────────────────────────────────

Subject Behind Token
    [Documentation]    The actor `/auth/introspect` resolves an access token to.
    [Arguments]    ${access_token}
    ${response}=    IDX Introspect With Token    ${access_token}
    Response Status Should Be    200    ${response}
    ${identity}=    Data Of    ${response}
    RETURN    ${identity}[subject]

Verified At
    [Documentation]    When the actor's email was verified, or None.
    [Arguments]    ${access_token}
    ${subject}=    Subject Behind Token    ${access_token}
    ${verified_at}=    Get From Dictionary    ${subject}    verified_at    default=${None}
    RETURN    ${verified_at}
