*** Settings ***
Documentation       Shared ground for every IdentityX step file: the scenario state,
...                 the stack, who a scenario is about, and the outcomes any "When"
...                 can end in.
...
...                 State lives in test-scoped variables: the scenario's actor
...                 (${EMAIL}/${PASSWORD}) and their ${TOKENS}, the ${RESPONSE} the
...                 last "When" got back for the "Then"s, and whatever tenant the
...                 "Given"s built (${ORGANIZATION}, ${PROJECT}, ...).
Library             Collections
Resource            ../../common.resource
Resource            ../../mailpit.resource
Resource            ../api.resource


*** Variables ***
${EMAIL}                ${None}
${PASSWORD}             ${None}
${RESPONSE}             ${None}
${TOKENS}               ${None}
${OLD_TOKENS}           ${None}
${ORGANIZATION}         ${None}
${PROJECT}              ${None}
# The project's owner, kept when the actor becomes one of its users.
${OWNER_TOKENS}         ${None}
# A second account: the one invited, added or looked up by the actor.
${OTHER_EMAIL}          ${None}
${OTHER_TOKENS}         ${None}

${VERIFY_SUBJECT}       Verify your email address
${RESET_SUBJECT}        Reset your password
${TOS_SUBJECT}          Our Terms of Service are changing


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

A visitor with the email "${email}"
    VAR    ${EMAIL}    ${email}    scope=TEST
    VAR    ${PASSWORD}    ${STRONG_PASSWORD}    scope=TEST

A registered visitor
    A visitor with a strong password
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}
    Response Status Should Be    201    ${response}

A verified user
    A registered visitor
    ${token}=    Action Token From Email    ${EMAIL}    ${VERIFY_SUBJECT}
    ${response}=    IDX Verify Email    ${token}
    Response Should Be Successful    ${response}

A signed-in platform user
    A registered visitor
    ${tokens}=    Tokens For    ${EMAIL}    ${PASSWORD}
    VAR    ${TOKENS}    ${tokens}    scope=TEST

The platform admin
    VAR    ${EMAIL}    ${ADMIN_EMAIL}    scope=TEST
    VAR    ${PASSWORD}    ${ADMIN_PASSWORD}    scope=TEST

A signed-in user who owns an organization
    A signed-in platform user
    ${slug}=    Unique Slug    org
    ${response}=    IDX Create Organization    ${TOKENS}[access_token]    E2E Labs    ${slug}
    Response Status Should Be    201    ${response}
    ${organization}=    Data Of    ${response}
    VAR    ${ORGANIZATION}    ${organization}    scope=TEST

A signed-in user who owns a project
    A signed-in user who owns an organization
    ${response}=    Create A Project In The Organization    E2E Checkout
    Response Status Should Be    201    ${response}
    ${project}=    Data Of    ${response}
    VAR    ${PROJECT}    ${project}    scope=TEST
    VAR    ${OWNER_TOKENS}    ${TOKENS}    scope=TEST

A user signed in to that project
    [Documentation]    Swaps the scenario's actor for a project user of ${PROJECT};
    ...    the owner stays reachable through ${OWNER_TOKENS}.
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    VAR    ${PASSWORD}    ${STRONG_PASSWORD}    scope=TEST
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    Response Status Should Be    201    ${response}
    ${tokens}=    Tokens For    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    VAR    ${TOKENS}    ${tokens}    scope=TEST

Another platform user
    [Documentation]    A second, signed-in platform account in ${OTHER_EMAIL}/${OTHER_TOKENS}.
    ${email}=    Unique Email
    ${response}=    IDX Register    ${email}    ${STRONG_PASSWORD}
    Response Status Should Be    201    ${response}
    ${tokens}=    Tokens For    ${email}    ${STRONG_PASSWORD}
    VAR    ${OTHER_EMAIL}    ${email}    scope=TEST
    VAR    ${OTHER_TOKENS}    ${tokens}    scope=TEST

Another user of that project
    [Documentation]    A second project user of ${PROJECT} in ${OTHER_EMAIL}/${OTHER_TOKENS}.
    ${email}=    Unique Email
    ${response}=    IDX Register    ${email}    ${STRONG_PASSWORD}    project_id=${PROJECT}[id]
    Response Status Should Be    201    ${response}
    ${tokens}=    Tokens For    ${email}    ${STRONG_PASSWORD}    project_id=${PROJECT}[id]
    VAR    ${OTHER_EMAIL}    ${email}    scope=TEST
    VAR    ${OTHER_TOKENS}    ${tokens}    scope=TEST

The other user takes over
    [Documentation]    Makes ${OTHER_EMAIL} the scenario's actor.
    VAR    ${EMAIL}    ${OTHER_EMAIL}    scope=TEST
    VAR    ${PASSWORD}    ${STRONG_PASSWORD}    scope=TEST
    VAR    ${TOKENS}    ${OTHER_TOKENS}    scope=TEST

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

Nothing is found
    Response Status Should Be    404    ${RESPONSE}

The error points at the "${field}" field
    ${error}=    Error Of    ${RESPONSE}
    Dictionary Should Contain Key    ${error}    fields    msg=No field errors: ${RESPONSE.text}
    ${names}=    Evaluate    [f["field"] for f in $error["fields"]]
    List Should Contain Value    ${names}    ${field}

The error says "${text}"
    ${error}=    Error Of    ${RESPONSE}
    Should Contain    ${error}[message]    ${text}    ignore_case=${True}

# ── helpers ───────────────────────────────────────────────────────────────

Tokens For
    [Documentation]    Signs in and returns the token pair, failing loudly otherwise.
    [Arguments]    ${email}    ${password}    ${project_id}=${None}
    ${response}=    IDX Login    ${email}    ${password}    project_id=${project_id}
    Response Status Should Be    200    ${response}
    ${tokens}=    Data Of    ${response}
    RETURN    ${tokens}

Create A Project In The Organization
    [Documentation]    A project needs a domain: its users' verify/reset links point
    ...    there. Domains are unique across projects, so each gets its own
    ...    *.localhost subdomain, which still resolves to the front.
    [Arguments]    ${name}    ${slug}=${None}
    IF    $slug is None
        ${slug}=    Unique Brand Slug    prj
    END
    ${domain}=    Unique Project Domain
    ${response}=    IDX Create Organization Project
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${name}    ${slug}    domain=${domain}
    RETURN    ${response}

Unique Project Domain
    ${label}=    Unique Brand Slug    dom
    ${domain}=    Evaluate    $IDX_APP_URL.replace("://", f"://{$label}.", 1)
    RETURN    ${domain}

Remember The Response
    [Arguments]    ${response}
    VAR    ${RESPONSE}    ${response}    scope=TEST

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

Ids Of
    [Documentation]    The `${key}` of every item in a list payload.
    [Arguments]    ${response}    ${key}=id
    ${items}=    Data Of    ${response}
    ${ids}=    Evaluate    [item[$key] for item in ($items or [])]
    RETURN    ${ids}
