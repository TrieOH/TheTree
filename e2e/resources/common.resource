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
# Browser suites open a visible Chromium by default so you can watch them;
# `just e2e ... --headless` (or no display at all) flips this.
${HEADLESS}             ${False}
# Playwright's per-action delay, e.g. 250ms, to follow a headed run by eye.
${SLOW_MO}              0s

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

Unique Brand Slug
    [Documentation]    A project brand slug: the DB only takes lowercase letters (^[a-z]+$, 3-32).
    [Arguments]    ${prefix}=prj
    ${id}=    Evaluate    ''.join(random.choices(string.ascii_lowercase, k=12))    modules=random,string
    RETURN    ${prefix}${id}

Is Envelope
    [Documentation]    True for a `fun.Response` envelope: a numeric `code` next to
    ...    `data` or `module` (no payload has those fields). The backends are moving
    ...    off it one at a time; drop this once none sends it.
    [Arguments]    ${body}
    ${enveloped}=    Evaluate    isinstance($body, dict) and isinstance($body.get("code"), int) and ("data" in $body or "module" in $body)
    RETURN    ${enveloped}

Data Of
    [Documentation]    The payload of a response: the body itself, or the `data`
    ...    of a `fun.Response` envelope.
    [Arguments]    ${response}
    ${body}=    Set Variable    ${response.json()}
    ${enveloped}=    Is Envelope    ${body}
    IF    not ${enveloped}    RETURN    ${body}
    Dictionary Should Contain Key    ${body}    data    msg=No data in response: ${response.text}
    RETURN    ${body}[data]

Items Of
    [Documentation]    The items of a list response: its `items`, or the bare array
    ...    a `fun.Response` envelope carries as data.
    [Arguments]    ${response}
    ${data}=    Data Of    ${response}
    ${items}=    Evaluate    $data if isinstance($data, list) else $data["items"]
    RETURN    ${items}

Error Of
    [Documentation]    The error of a failed response as `{code, message, fields}`:
    ...    the body itself, an RFC 9457 problem's `code`, `detail` and `errors`, or
    ...    the `error` object of a `fun.Response` envelope.
    [Arguments]    ${response}
    ${body}=    Set Variable    ${response.json()}
    ${enveloped}=    Is Envelope    ${body}
    IF    not ${enveloped}
        ${error}=    Evaluate    {k: v for k, v in {"code": $body.get("code"), "message": $body.get("message") or $body.get("detail") or $body.get("title"), "fields": $body.get("fields") or $body.get("errors")}.items() if v is not None}
        RETURN    ${error}
    END
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
