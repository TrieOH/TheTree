*** Settings ***
Documentation       Shared plumbing for every suite: where the e2e stack lives,
...                 unique test data, and helpers for the `fun.Response`
...                 envelope every backend wraps its JSON in.
Library             Collections
Library             RequestsLibrary


*** Variables ***
# Ports match e2e/compose.yml and `just e2e-up`. Override with
# `robot --variable IDX_URL:...` to point a run at another stack.
${IDX_URL}              http://localhost:18080
${IDX_APP_URL}          http://localhost:13001
${MAILPIT_URL}          http://localhost:18025
${HEADLESS}             ${True}

# The platform super admin `/auth/setup` creates once per stack.
${ADMIN_EMAIL}          admin@e2e.trieoh.test
${ADMIN_PASSWORD}       E2e-Admin-Passw0rd!

# Meets the passwd rule: upper, number, symbol, 8–72 chars.
${STRONG_PASSWORD}      E2e-Passw0rd!


*** Keywords ***
Unique Email
    [Documentation]    A fresh address per call, so scenarios never share accounts
    ...    and the suite can rerun against the same stack.
    ${id}=    Evaluate    uuid.uuid4().hex[:12]    modules=uuid
    RETURN    e2e-${id}@e2e.trieoh.test

Unique Slug
    [Arguments]    ${prefix}=e2e
    ${id}=    Evaluate    uuid.uuid4().hex[:8]    modules=uuid
    RETURN    ${prefix}-${id}

Data Of
    [Documentation]    The `data` payload of a `fun.Response` envelope.
    [Arguments]    ${response}
    ${body}=    Set Variable    ${response.json()}
    Dictionary Should Contain Key    ${body}    data    msg=No data in response: ${response.text}
    RETURN    ${body}[data]

Error Of
    [Documentation]    The `error` object of a failed `fun.Response` envelope.
    [Arguments]    ${response}
    ${body}=    Set Variable    ${response.json()}
    Dictionary Should Contain Key    ${body}    error    msg=No error in response: ${response.text}
    RETURN    ${body}[error]

Response Status Should Be
    [Documentation]    Like RequestsLibrary's `Status Should Be`, but shows the body on a miss.
    [Arguments]    ${expected}    ${response}
    Should Be Equal As Integers    ${response.status_code}    ${expected}
    ...    msg=Expected HTTP ${expected}, got ${response.status_code}: ${response.text}

Response Should Be Successful
    [Arguments]    ${response}
    Should Be True    200 <= ${response.status_code} < 300
    ...    msg=Expected a 2xx, got ${response.status_code}: ${response.text}

Bearer
    [Documentation]    Headers authenticating as the holder of an access token.
    [Arguments]    ${access_token}
    &{headers}=    Create Dictionary    Authorization=Bearer ${access_token}
    RETURN    ${headers}

Wait Until Service Is Healthy
    [Arguments]    ${base_url}    ${timeout}=60s
    Wait Until Keyword Succeeds    ${timeout}    1s    Service Answers Health    ${base_url}

Service Answers Health
    [Arguments]    ${base_url}
    ${response}=    GET    ${base_url}/health    expected_status=200
    RETURN    ${response}
