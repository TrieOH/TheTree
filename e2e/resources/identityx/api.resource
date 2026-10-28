*** Settings ***
Documentation       Raw IdentityX HTTP calls, one keyword per operation in
...                 api/identityx/api-spec.yml. They return the response and
...                 never assert on its status: what counts as success is the
...                 scenario's call, made in the step files.
Library             Collections
Library             RequestsLibrary
Resource            ../common.resource


*** Keywords ***
Open IdentityX Session
    Create Session    idx    ${IDX_URL}    disable_warnings=1

# ── system ────────────────────────────────────────────────────────────────

IDX Get Health
    ${response}=    GET On Session    idx    /health    expected_status=any
    RETURN    ${response}

IDX Get OpenAPI Spec
    ${response}=    GET On Session    idx    /docs/openapi.yml    expected_status=any
    RETURN    ${response}

IDX Get Setup Status
    ${response}=    GET On Session    idx    /auth/setup    expected_status=any
    RETURN    ${response}

IDX Set Up Platform
    [Arguments]    ${email}    ${password}
    &{body}=    Create Dictionary    email=${email}    password=${password}
    ${response}=    POST On Session    idx    /auth/setup    json=${body}    expected_status=any
    RETURN    ${response}

# ── authn ─────────────────────────────────────────────────────────────────

IDX Register
    [Arguments]    ${email}    ${password}    ${project_id}=${None}    ${accepted_tos}=${None}
    &{body}=    Create Dictionary    email=${email}    password=${password}
    IF    $accepted_tos is not None
        Set To Dictionary    ${body}    accepted_tos=${accepted_tos}
    END
    ${params}=    Project Query    ${project_id}
    ${response}=    POST On Session    idx    /auth/register
    ...    json=${body}    params=${params}    expected_status=any
    RETURN    ${response}

IDX Login
    [Arguments]    ${email}    ${password}    ${project_id}=${None}
    &{body}=    Create Dictionary    email=${email}    password=${password}
    ${params}=    Project Query    ${project_id}
    ${response}=    POST On Session    idx    /auth/login
    ...    json=${body}    params=${params}    expected_status=any
    RETURN    ${response}

IDX Refresh
    [Arguments]    ${refresh_token}
    &{headers}=    Create Dictionary    Refresh-Token=${refresh_token}
    ${response}=    POST On Session    idx    /auth/refresh    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Logout
    [Arguments]    ${access_token}    ${refresh_token}
    ${headers}=    Bearer    ${access_token}
    Set To Dictionary    ${headers}    Refresh-Token=${refresh_token}
    ${response}=    POST On Session    idx    /auth/logout    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Introspect With Token
    [Arguments]    ${access_token}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /auth/introspect    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Introspect With API Key
    [Arguments]    ${raw_key}
    &{headers}=    Create Dictionary    X-API-KEY=${raw_key}
    ${response}=    GET On Session    idx    /auth/introspect    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Introspect Anonymously
    ${response}=    GET On Session    idx    /auth/introspect    expected_status=any
    RETURN    ${response}

IDX Verify Email
    [Arguments]    ${token}
    &{body}=    Create Dictionary    token=${token}
    ${response}=    POST On Session    idx    /auth/verify-email    json=${body}    expected_status=any
    RETURN    ${response}

IDX Resend Verification
    [Arguments]    ${email}    ${project_id}=${None}
    &{body}=    Create Dictionary    email=${email}
    IF    $project_id is not None
        Set To Dictionary    ${body}    project_id=${project_id}
    END
    ${response}=    POST On Session    idx    /auth/resend-verification    json=${body}    expected_status=any
    RETURN    ${response}

IDX Forgot Password
    [Arguments]    ${email}    ${project_id}=${None}
    &{body}=    Create Dictionary    email=${email}
    IF    $project_id is not None
        Set To Dictionary    ${body}    project_id=${project_id}
    END
    ${response}=    POST On Session    idx    /auth/forgot-password    json=${body}    expected_status=any
    RETURN    ${response}

IDX Reset Password
    [Arguments]    ${token}    ${password}
    &{body}=    Create Dictionary    token=${token}    password=${password}
    ${response}=    POST On Session    idx    /auth/reset-password    json=${body}    expected_status=any
    RETURN    ${response}

IDX Get JWKS
    [Arguments]    ${project_id}=${None}
    ${params}=    Project Query    ${project_id}
    ${response}=    GET On Session    idx    /.well-known/jwks.json    params=${params}    expected_status=any
    RETURN    ${response}

# ── social login ──────────────────────────────────────────────────────────

IDX Discover OAuth Providers
    [Arguments]    ${project_id}=${None}
    ${params}=    Project Query    ${project_id}
    ${response}=    GET On Session    idx    /auth/oauth-providers    params=${params}    expected_status=any
    RETURN    ${response}

IDX Start Social Login
    [Arguments]    ${provider}    ${project_id}=${None}
    ${params}=    Project Query    ${project_id}
    ${response}=    GET On Session    idx    /auth/${provider}/connect    params=${params}    expected_status=any
    RETURN    ${response}

IDX List Project OAuth Providers
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/oauth-providers
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Configure OAuth Provider
    [Arguments]    ${access_token}    ${project_id}    ${provider}    ${client_id}    ${client_secret}    ${callback_url}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    provider=${provider}    client_id=${client_id}
    ...    client_secret=${client_secret}    callback_url=${callback_url}
    ${response}=    POST On Session    idx    /projects/${project_id}/oauth-providers
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Update OAuth Provider
    [Arguments]    ${access_token}    ${provider_id}    &{fields}
    ${headers}=    Bearer    ${access_token}
    ${response}=    PATCH On Session    idx    /oauth-providers/${provider_id}
    ...    json=${fields}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Disable OAuth Provider
    [Arguments]    ${access_token}    ${provider_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    POST On Session    idx    /oauth-providers/${provider_id}/disable
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Enable OAuth Provider
    [Arguments]    ${access_token}    ${provider_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    POST On Session    idx    /oauth-providers/${provider_id}/enable
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Delete OAuth Provider
    [Arguments]    ${access_token}    ${provider_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    DELETE On Session    idx    /oauth-providers/${provider_id}
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── organizations ─────────────────────────────────────────────────────────

IDX Create Organization
    [Arguments]    ${access_token}    ${name}    ${slug}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    name=${name}    slug=${slug}
    ${response}=    POST On Session    idx    /organizations
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Organizations
    [Arguments]    ${access_token}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /organizations    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Organization Members
    [Arguments]    ${access_token}    ${organization_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /organizations/${organization_id}/members
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Add Organization Member
    [Arguments]    ${access_token}    ${organization_id}    ${email}    ${role}=member
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    actor_email=${email}    role=${role}
    ${response}=    POST On Session    idx    /organizations/${organization_id}/members
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Remove Organization Member
    [Arguments]    ${access_token}    ${organization_id}    ${email}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    actor_email=${email}
    ${response}=    DELETE On Session    idx    /organizations/${organization_id}/members
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Organization Projects
    [Arguments]    ${access_token}    ${organization_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /organizations/${organization_id}/projects
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Create Organization Project
    [Arguments]    ${access_token}    ${organization_id}    ${name}    ${brand_slug}    ${domain}=${None}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    name=${name}    brand_slug=${brand_slug}
    IF    $domain is not None
        Set To Dictionary    ${body}    domain=${domain}
    END
    ${response}=    POST On Session    idx    /organizations/${organization_id}/projects
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Organization Project Members
    [Arguments]    ${access_token}    ${organization_id}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /organizations/${organization_id}/projects/${project_id}/members
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── projects, members, actors ─────────────────────────────────────────────

IDX List Project Members
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/members
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Add Project Member
    [Arguments]    ${access_token}    ${project_id}    ${email}    ${role}=member
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    actor_email=${email}    role=${role}
    ${response}=    POST On Session    idx    /projects/${project_id}/members
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Remove Project Member
    [Arguments]    ${access_token}    ${project_id}    ${email}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    actor_email=${email}
    ${response}=    DELETE On Session    idx    /projects/${project_id}/members
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Actors
    [Arguments]    ${headers}    ${project_id}
    ${response}=    GET On Session    idx    /projects/${project_id}/actors
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Get Actor
    [Arguments]    ${access_token}    ${project_id}    ${actor_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/actors/${actor_id}
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Get Actor By Email
    [Arguments]    ${access_token}    ${project_id}    ${email}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/actors/${email}:by_email
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Create Actor
    [Arguments]    ${access_token}    ${project_id}    ${auth_method}    ${type}    ${email}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    auth_method=${auth_method}    type=${type}    email=${email}
    ${response}=    POST On Session    idx    /projects/${project_id}/actors
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── capabilities, api keys ────────────────────────────────────────────────

IDX List Capabilities
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/capabilities
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Create Capability
    [Arguments]    ${access_token}    ${project_id}    ${resource}    ${action}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    resource=${resource}    action=${action}
    ${response}=    POST On Session    idx    /projects/${project_id}/capabilities
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Create API Key
    [Documentation]    Bound to the caller (no subject_id). Optionally restricted
    ...    to capability ids and given an expiry.
    [Arguments]    ${access_token}    ${project_id}    ${name}    ${capabilities}=${None}    ${expires_at}=${None}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    name=${name}    env=sandbox
    IF    $capabilities is not None
        Set To Dictionary    ${body}    capabilities=${capabilities}
    END
    IF    $expires_at is not None
        Set To Dictionary    ${body}    expires_at=${expires_at}
    END
    ${response}=    POST On Session    idx    /projects/${project_id}/api_keys
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── email templates ───────────────────────────────────────────────────────

IDX List Email Templates
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/email-templates
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Get Email Template
    [Arguments]    ${access_token}    ${project_id}    ${kind}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/email-templates/${kind}
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Put Email Template
    [Arguments]    ${access_token}    ${project_id}    ${kind}    ${subject}    ${body}
    ${headers}=    Bearer    ${access_token}
    &{payload}=    Create Dictionary    subject=${subject}    body=${body}
    ${response}=    PUT On Session    idx    /projects/${project_id}/email-templates/${kind}
    ...    json=${payload}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Delete Email Template
    [Arguments]    ${access_token}    ${project_id}    ${kind}
    ${headers}=    Bearer    ${access_token}
    ${response}=    DELETE On Session    idx    /projects/${project_id}/email-templates/${kind}
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── terms of service ──────────────────────────────────────────────────────

IDX Get Terms
    [Arguments]    ${project_id}
    ${response}=    GET On Session    idx    /projects/${project_id}/terms-of-service    expected_status=any
    RETURN    ${response}

IDX Introduce Terms
    [Arguments]    ${access_token}    ${project_id}    ${content}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    content=${content}
    ${response}=    POST On Session    idx    /projects/${project_id}/terms-of-service
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Update Terms
    [Arguments]    ${access_token}    ${project_id}    ${content}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    content=${content}
    ${response}=    PUT On Session    idx    /projects/${project_id}/terms-of-service
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX List Terms Acceptances
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/terms-of-service/acceptances
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Accept Terms
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    POST On Session    idx    /projects/${project_id}/terms-of-service/acceptances
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── profiles ──────────────────────────────────────────────────────────────

IDX Get Project Profile Schema
    [Arguments]    ${project_id}
    ${response}=    GET On Session    idx    /projects/${project_id}/profile-schema    expected_status=any
    RETURN    ${response}

IDX Put Project Profile Schema
    [Arguments]    ${access_token}    ${project_id}    ${schema}    ${active}=${True}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    schema=${schema}    active=${active}
    ${response}=    PUT On Session    idx    /projects/${project_id}/profile-schema
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Get Project Profile
    [Arguments]    ${project_id}    ${actor_id}
    ${response}=    GET On Session    idx    /projects/${project_id}/actors/${actor_id}/profile    expected_status=any
    RETURN    ${response}

IDX Put Project Profile
    [Documentation]    ${handle} is only sent when given, so the call can leave it untouched.
    [Arguments]    ${access_token}    ${project_id}    ${actor_id}    ${profile}    ${handle}=${None}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    profile=${profile}
    IF    $handle is not None
        Set To Dictionary    ${body}    handle=${handle}
    END
    ${response}=    PUT On Session    idx    /projects/${project_id}/actors/${actor_id}/profile
    ...    json=${body}    headers=${headers}    expected_status=any
    RETURN    ${response}

IDX Get Profile By Handle
    [Arguments]    ${handle}
    ${response}=    GET On Session    idx    /profiles/by-handle/${handle}    expected_status=any
    RETURN    ${response}

IDX List Outdated Project Profiles
    [Arguments]    ${access_token}    ${project_id}
    ${headers}=    Bearer    ${access_token}
    ${response}=    GET On Session    idx    /projects/${project_id}/profiles/outdated
    ...    headers=${headers}    expected_status=any
    RETURN    ${response}

# ── helpers ───────────────────────────────────────────────────────────────

Project Query
    [Documentation]    `?project_id=` for project-scoped calls, nothing for platform ones.
    [Arguments]    ${project_id}
    &{params}=    Create Dictionary
    IF    $project_id is not None
        Set To Dictionary    ${params}    project_id=${project_id}
    END
    RETURN    ${params}
