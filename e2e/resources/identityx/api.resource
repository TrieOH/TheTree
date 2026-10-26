*** Settings ***
Documentation       Raw IdentityX HTTP calls, one keyword per operation in
...                 api/identityx/api-spec.yml. They return the response and
...                 never assert on its status: what counts as success is the
...                 scenario's call, made in steps.resource.
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

IDX Verify Email
    [Arguments]    ${token}
    &{body}=    Create Dictionary    token=${token}
    ${response}=    POST On Session    idx    /auth/verify-email    json=${body}    expected_status=any
    RETURN    ${response}

IDX Forgot Password
    [Arguments]    ${email}
    &{body}=    Create Dictionary    email=${email}
    ${response}=    POST On Session    idx    /auth/forgot-password    json=${body}    expected_status=any
    RETURN    ${response}

IDX Reset Password
    [Arguments]    ${token}    ${password}
    &{body}=    Create Dictionary    token=${token}    password=${password}
    ${response}=    POST On Session    idx    /auth/reset-password    json=${body}    expected_status=any
    RETURN    ${response}

# ── organizations, projects, api keys ─────────────────────────────────────

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

IDX Create API Key
    [Documentation]    Bound to the caller (no subject_id) and unrestricted.
    [Arguments]    ${access_token}    ${project_id}    ${name}
    ${headers}=    Bearer    ${access_token}
    &{body}=    Create Dictionary    name=${name}    env=sandbox
    ${response}=    POST On Session    idx    /projects/${project_id}/api_keys
    ...    json=${body}    headers=${headers}    expected_status=any
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
