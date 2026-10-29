*** Settings ***
Documentation       Sentences about what a project owner configures over the API:
...                 capabilities and API keys, email templates, terms of service,
...                 profile schemas and profiles, and social login providers.
Library             Collections
Resource            common.resource


*** Variables ***
${CAPABILITY}           ${None}
${API_KEY}              ${None}
${OAUTH_PROVIDER}       ${None}
${HANDLE}               ${None}
${CUSTOM_SUBJECT}       Welcome to {{.ProjectName}}: confirm your email
${CUSTOM_BODY}          <p>Hi {{.Email}}!</p><p><a href="{{.ActionURL}}">Confirm</a></p>
# `full_name` is required, so an empty profile fails the schema.
&{PROFILE_SCHEMA}       type=object    required=${{["full_name"]}}
...                     properties=${{{"full_name": {"type": "string"}}}}


*** Keywords ***
# ── capabilities and api keys ─────────────────────────────────────────────

They define the capability "${resource}:${action}"
    ${response}=    IDX Create Capability    ${TOKENS}[access_token]    ${PROJECT}[id]    ${resource}    ${action}
    Remember The Response    ${response}

The project has the capability "${resource}:${action}"
    ${response}=    IDX Create Capability    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${resource}    ${action}
    Response Status Should Be    201    ${response}
    ${capability}=    Data Of    ${response}
    VAR    ${CAPABILITY}    ${capability}    scope=TEST

The capability is created
    Response Status Should Be    201    ${RESPONSE}
    ${capability}=    Data Of    ${RESPONSE}
    VAR    ${CAPABILITY}    ${capability}    scope=TEST

The capability shows in the project's capabilities
    ${response}=    IDX List Capabilities    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${ids}=    Ids Of    ${response}
    List Should Contain Value    ${ids}    ${CAPABILITY}[id]

They create an API key for the project
    ${response}=    IDX Create API Key    ${TOKENS}[access_token]    ${PROJECT}[id]    e2e
    Remember The Response    ${response}

They create an API key restricted to that capability
    ${capabilities}=    Create List    ${CAPABILITY}[id]
    ${response}=    IDX Create API Key    ${TOKENS}[access_token]    ${PROJECT}[id]    e2e-restricted
    ...    capabilities=${capabilities}
    Remember The Response    ${response}

They create an API key that already expired
    ${response}=    IDX Create API Key    ${TOKENS}[access_token]    ${PROJECT}[id]    e2e-expired
    ...    expires_at=2020-01-01T00:00:00Z
    Remember The Response    ${response}

Someone authenticates with a made-up API key
    ${response}=    IDX Introspect With API Key    idx_v1_sandbox_not-a-real-key
    Remember The Response    ${response}

The raw API key is handed over
    Response Status Should Be    201    ${RESPONSE}
    ${created}=    Data Of    ${RESPONSE}
    Should Not Be Empty    ${created}[raw_key]
    VAR    ${API_KEY}    ${created}[raw_key]    scope=TEST

Only the raw key and public details are returned
    [Documentation]    The key record is for display: the hash that verifies the
    ...    key is server-side state and has no business on the wire.
    ${created}=    Data Of    ${RESPONSE}
    Dictionary Should Not Contain Key    ${created}[key]    key_hash
    ...    msg=The API key response leaks key_hash: ${created}[key]

The API key authenticates as them
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Introspect With API Key    ${API_KEY}
    Response Status Should Be    200    ${response}
    ${identity}=    Data Of    ${response}
    Should Be Equal    ${identity}[subject][id]    ${me}[id]
    Should Be Equal    ${identity}[cred][type]    api_key

The API key carries the capability "${grant}"
    ${response}=    IDX Introspect With API Key    ${API_KEY}
    Response Status Should Be    200    ${response}
    ${identity}=    Data Of    ${response}
    List Should Contain Value    ${identity}[subject][capabilities]    ${grant}

The API key is refused
    ${response}=    IDX Introspect With API Key    ${API_KEY}
    Response Status Should Be    401    ${response}

The API key can list the project's actors
    &{headers}=    Create Dictionary    X-API-KEY=${API_KEY}
    ${response}=    IDX List Actors    ${headers}    ${PROJECT}[id]
    Response Status Should Be    200    ${response}

# ── email templates ───────────────────────────────────────────────────────

They list the project's email templates
    ${response}=    IDX List Email Templates    ${TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

They customize the "${kind}" email
    ${response}=    IDX Put Email Template    ${TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    ...    ${CUSTOM_SUBJECT}    ${CUSTOM_BODY}
    Remember The Response    ${response}

The project has a customized "${kind}" email
    ${response}=    IDX Put Email Template    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    ...    ${CUSTOM_SUBJECT}    ${CUSTOM_BODY}
    Response Status Should Be    200    ${response}

They save a "verify" email without the action link
    ${response}=    IDX Put Email Template    ${TOKENS}[access_token]    ${PROJECT}[id]    verify
    ...    Hello    <p>No link in here</p>
    Remember The Response    ${response}

They save a "verify" email that does not render
    ${response}=    IDX Put Email Template    ${TOKENS}[access_token]    ${PROJECT}[id]    verify
    ...    Hello    <a href="{{.ActionURL}}">go</a>{{.NoSuchField
    Remember The Response    ${response}

They ask for the "${kind}" email template
    ${response}=    IDX Get Email Template    ${TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    Remember The Response    ${response}

They reset the "${kind}" email to the default
    ${response}=    IDX Delete Email Template    ${TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    Remember The Response    ${response}

Every kind of email comes with the default template
    Response Status Should Be    200    ${RESPONSE}
    ${templates}=    Items Of    ${RESPONSE}
    ${sources}=    Evaluate    {t["kind"]: t["source"] for t in $templates}
    &{expected}=    Create Dictionary    verify=default    reset=default    tos=default
    Dictionaries Should Be Equal    ${sources}    ${expected}

The "${kind}" email is the project's own
    ${response}=    IDX Get Email Template    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    Response Status Should Be    200    ${response}
    ${template}=    Data Of    ${response}
    Should Be Equal    ${template}[source]    override
    Should Be Equal    ${template}[subject]    ${CUSTOM_SUBJECT}

The "${kind}" email is back to the default
    ${response}=    IDX Get Email Template    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${kind}
    Response Status Should Be    200    ${response}
    ${template}=    Data Of    ${response}
    Should Be Equal    ${template}[source]    default

The verification email they get uses the project's template
    ${subject}=    Evaluate    $CUSTOM_SUBJECT.replace("{{.ProjectName}}", $PROJECT["name"])
    ${message_id}=    Wait For Email To    ${EMAIL}    ${subject}
    ${body}=    Body Of Email    ${message_id}
    Should Contain    ${body}    Hi ${EMAIL}!
    Should Contain    ${body}    token=

# ── terms of service ──────────────────────────────────────────────────────

Someone reads the project's terms
    ${response}=    IDX Get Terms    ${PROJECT}[id]
    Remember The Response    ${response}

They introduce the terms "${content}"
    ${response}=    IDX Introduce Terms    ${TOKENS}[access_token]    ${PROJECT}[id]    ${content}
    Remember The Response    ${response}

They update the terms to "${content}"
    ${response}=    IDX Update Terms    ${TOKENS}[access_token]    ${PROJECT}[id]    ${content}
    Remember The Response    ${response}

The project has the terms "${content}"
    ${response}=    IDX Introduce Terms    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${content}
    Response Status Should Be    200    ${response}

They sign up to the project accepting the terms
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]    accepted_tos=${True}
    Remember The Response    ${response}

They sign up to the project without accepting the terms
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    Remember The Response    ${response}

They sign up to the project declining the terms
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]    accepted_tos=${False}
    Remember The Response    ${response}

They accept the current terms
    ${response}=    IDX Accept Terms    ${TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

They accept the current terms again
    They accept the current terms

They read the terms acceptances
    ${response}=    IDX List Terms Acceptances    ${TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

The terms are at version ${version}
    Response Status Should Be    200    ${RESPONSE}
    ${tos}=    Data Of    ${RESPONSE}
    Should Be Equal As Integers    ${tos}[version]    ${version}

The terms take effect 30 days after they are announced
    ${tos}=    Data Of    ${RESPONSE}
    ${days}=    Evaluate
    ...    (datetime.datetime.fromisoformat($tos["effective_at"]) - datetime.datetime.fromisoformat($tos["updated_at"])).days
    ...    modules=datetime
    Should Be True    29 <= ${days} <= 30    msg=Terms take effect ${days} days out

The project's terms read "${content}"
    ${response}=    IDX Get Terms    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${tos}=    Data Of    ${response}
    Should Be Equal    ${tos}[content]    ${content}

They are told about the new terms by email
    Wait For Email To    ${EMAIL}    ${TOS_SUBJECT}

Their acceptance of version ${version} is on record
    ${response}=    IDX List Terms Acceptances    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${acceptances}=    Items Of    ${response}
    ${mine}=    Evaluate    [a for a in $acceptances if a.get("actor_email") == $EMAIL]
    Should Not Be Empty    ${mine}    msg=No acceptance on record for ${EMAIL}: ${acceptances}
    Should Be Equal As Integers    ${mine}[0][tos_version]    ${version}
    Should Be Equal    ${mine}[0][source]    clickwrap
    Length Should Be    ${mine}[0][content_hash]    64

They have accepted the terms only once
    ${response}=    IDX List Terms Acceptances    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${acceptances}=    Items Of    ${response}
    ${mine}=    Evaluate    [a for a in $acceptances if a.get("actor_email") == $EMAIL]
    Length Should Be    ${mine}    1

# ── profiles ──────────────────────────────────────────────────────────────

Someone reads the project's profile schema
    ${response}=    IDX Get Project Profile Schema    ${PROJECT}[id]
    Remember The Response    ${response}

They publish a profile schema requiring a full name
    ${response}=    IDX Put Project Profile Schema    ${TOKENS}[access_token]    ${PROJECT}[id]    ${PROFILE_SCHEMA}
    Remember The Response    ${response}

The project requires a full name on profiles
    ${response}=    IDX Put Project Profile Schema    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${PROFILE_SCHEMA}
    Response Status Should Be    200    ${response}

The owner starts requiring an age as well
    &{properties}=    Create Dictionary    full_name=${{{"type": "string"}}}    age=${{{"type": "integer"}}}
    &{schema}=    Create Dictionary    type=object    required=${{["full_name", "age"]}}    properties=${properties}
    ${response}=    IDX Put Project Profile Schema    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${schema}
    Response Status Should Be    200    ${response}

They set their full name to "${name}" with a handle
    ${handle}=    Unique Brand Slug    h
    VAR    ${HANDLE}    ${handle}    scope=TEST
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=${name}
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    ...    handle=${handle}
    Remember The Response    ${response}

They have the full name "${name}" and a handle
    They set their full name to "${name}" with a handle
    Response Status Should Be    200    ${RESPONSE}

They change their full name to "${name}"
    [Documentation]    Sends no handle, which should leave the handle as it is.
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=${name}
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    Remember The Response    ${response}

They save a profile without a full name
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    Remember The Response    ${response}

They save a profile whose full name is a number
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=${42}
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    Remember The Response    ${response}

They try to edit the other user's profile
    ${other}=    Subject Behind Token    ${OTHER_TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=Not Yours
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${other}[id]    ${profile}
    Remember The Response    ${response}

They try to take that handle
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=Copycat
    ${response}=    IDX Put Project Profile    ${TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    ...    handle=${HANDLE}
    Remember The Response    ${response}

The owner renames the user to "${name}"
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    &{profile}=    Create Dictionary    full_name=${name}
    ${response}=    IDX Put Project Profile    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]    ${profile}
    Remember The Response    ${response}

Someone looks the profile up by its handle
    ${response}=    IDX Get Profile By Handle    ${HANDLE}
    Remember The Response    ${response}

Someone looks up a handle nobody has
    ${handle}=    Unique Brand Slug    nobody
    ${response}=    IDX Get Profile By Handle    ${handle}
    Remember The Response    ${response}

They list the outdated profiles
    ${response}=    IDX List Outdated Project Profiles    ${TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

The owner lists the outdated profiles
    ${response}=    IDX List Outdated Project Profiles    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

The profile schema is at version ${version}
    Response Status Should Be    200    ${RESPONSE}
    ${schema}=    Data Of    ${RESPONSE}
    Should Be Equal As Integers    ${schema}[version]    ${version}
    Should Be True    ${schema}[active]

Anyone can read the schema
    ${response}=    IDX Get Project Profile Schema    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${schema}=    Data Of    ${response}
    Should Be Equal    ${schema}[schema]    ${PROFILE_SCHEMA}

Their public profile shows the full name "${name}"
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Get Project Profile    ${PROJECT}[id]    ${me}[id]
    Response Status Should Be    200    ${response}
    ${profile}=    Data Of    ${response}
    Should Be Equal    ${profile}[profile][full_name]    ${name}

Their profile keeps its handle
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Get Project Profile    ${PROJECT}[id]    ${me}[id]
    Response Status Should Be    200    ${response}
    ${profile}=    Data Of    ${response}
    Should Be Equal    ${profile}[handle]    ${HANDLE}    msg=Saving without a handle changed it

The profile found is theirs
    Response Status Should Be    200    ${RESPONSE}
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${profile}=    Data Of    ${RESPONSE}
    Should Be Equal    ${profile}[actor_id]    ${me}[id]

Their profile is flagged as outdated
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Get Project Profile    ${PROJECT}[id]    ${me}[id]
    ${profile}=    Data Of    ${response}
    Should Be True    ${profile}[outdated]    msg=Profile is not flagged outdated: ${profile}

Their profile is in the list
    Response Status Should Be    200    ${RESPONSE}
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${ids}=    Ids Of    ${RESPONSE}    actor_id
    List Should Contain Value    ${ids}    ${me}[id]

# ── social login providers ────────────────────────────────────────────────

They configure GitHub login for the project
    ${callback}=    Set Variable    ${PROJECT}[domain]/auth/github/callback
    ${response}=    IDX Configure OAuth Provider    ${TOKENS}[access_token]    ${PROJECT}[id]
    ...    github    e2e-client-id    e2e-client-secret    ${callback}
    Remember The Response    ${response}

The project offers GitHub login
    ${actor}=    Set Variable If    $OWNER_TOKENS is not None    ${OWNER_TOKENS}    ${TOKENS}
    ${callback}=    Set Variable    ${PROJECT}[domain]/auth/github/callback
    ${response}=    IDX Configure OAuth Provider    ${actor}[access_token]    ${PROJECT}[id]
    ...    github    e2e-client-id    e2e-client-secret    ${callback}
    Response Status Should Be    201    ${response}
    ${provider}=    Data Of    ${response}
    VAR    ${OAUTH_PROVIDER}    ${provider}    scope=TEST

They disable GitHub login
    ${response}=    IDX Disable OAuth Provider    ${TOKENS}[access_token]    ${OAUTH_PROVIDER}[id]
    Remember The Response    ${response}

They enable GitHub login again
    ${response}=    IDX Enable OAuth Provider    ${TOKENS}[access_token]    ${OAUTH_PROVIDER}[id]
    Remember The Response    ${response}

They change the GitHub client id to "${client_id}"
    ${response}=    IDX Update OAuth Provider    ${TOKENS}[access_token]    ${OAUTH_PROVIDER}[id]
    ...    client_id=${client_id}
    Remember The Response    ${response}

They remove GitHub login
    ${response}=    IDX Delete OAuth Provider    ${TOKENS}[access_token]    ${OAUTH_PROVIDER}[id]
    Remember The Response    ${response}

They list the project's login providers
    ${response}=    IDX List Project OAuth Providers    ${TOKENS}[access_token]    ${PROJECT}[id]
    Remember The Response    ${response}

A visitor starts a GitHub login on the project
    ${response}=    IDX Start Social Login    github    project_id=${PROJECT}[id]
    Remember The Response    ${response}

A visitor starts a GitHub login on the platform
    ${response}=    IDX Start Social Login    github
    Remember The Response    ${response}

Someone asks which logins the platform offers
    ${response}=    IDX Discover OAuth Providers
    Remember The Response    ${response}

Someone asks which logins a project that does not exist offers
    ${id}=    Evaluate    str(uuid.uuid4())    modules=uuid
    ${response}=    IDX Discover OAuth Providers    project_id=${id}
    Remember The Response    ${response}

The provider is configured without exposing the secret
    Response Status Should Be    201    ${RESPONSE}
    ${provider}=    Data Of    ${RESPONSE}
    Should Be Equal    ${provider}[provider]    github
    Should Be True    ${provider}[enabled]
    Dictionary Should Not Contain Key    ${provider}    client_secret
    VAR    ${OAUTH_PROVIDER}    ${provider}    scope=TEST

Visitors are offered GitHub login
    ${response}=    IDX Discover OAuth Providers    project_id=${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${providers}=    Ids Of    ${response}    provider
    List Should Contain Value    ${providers}    github

Visitors are offered no social login
    ${response}=    IDX Discover OAuth Providers    project_id=${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${providers}=    Ids Of    ${response}    provider
    Should Be Empty    ${providers}

No social login is offered
    Response Status Should Be    200    ${RESPONSE}
    ${providers}=    Ids Of    ${RESPONSE}    provider
    Should Be Empty    ${providers}

They are sent to GitHub with client id "${client_id}"
    Response Status Should Be    200    ${RESPONSE}
    ${data}=    Data Of    ${RESPONSE}
    Should Start With    ${data}[url]    https://github.com/login/oauth/authorize?
    Should Contain    ${data}[url]    client_id=${client_id}
    Should Contain    ${data}[url]    state=
