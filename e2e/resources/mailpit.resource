*** Settings ***
Documentation       Reading the e2e stack's inbox through mailpit's HTTP API
...                 (https://mailpit.axllent.org/docs/api-v1/). Emails are sent
...                 by background jobs, so every read waits for the message.
Library             Collections
Library             String
Library             RequestsLibrary
Resource            common.resource


*** Variables ***
${EMAIL_TIMEOUT}        30s


*** Keywords ***
Open Mailpit Session
    Create Session    mailpit    ${MAILPIT_URL}    disable_warnings=1

Emails To
    [Documentation]    Messages sent to ${address} with ${subject}, newest first.
    [Arguments]    ${address}    ${subject}
    &{params}=    Create Dictionary    query=to:"${address}" subject:"${subject}"
    ${response}=    GET On Session    mailpit    /api/v1/search    params=${params}
    RETURN    ${response.json()}[messages]

Newest Email To
    [Arguments]    ${address}    ${subject}
    ${messages}=    Emails To    ${address}    ${subject}
    Should Not Be Empty    ${messages}    msg=No "${subject}" email to ${address} yet
    RETURN    ${messages}[0]

Wait For Email To
    [Documentation]    Waits for at least ${count} "${subject}" emails to ${address}
    ...    and returns the newest one's id.
    [Arguments]    ${address}    ${subject}    ${count}=1
    Wait Until Keyword Succeeds    ${EMAIL_TIMEOUT}    1s
    ...    Email Count Should Be At Least    ${address}    ${subject}    ${count}
    ${message}=    Newest Email To    ${address}    ${subject}
    RETURN    ${message}[ID]

Email Count Should Be At Least
    [Arguments]    ${address}    ${subject}    ${count}
    ${messages}=    Emails To    ${address}    ${subject}
    ${found}=    Get Length    ${messages}
    Should Be True    ${found} >= ${count}
    ...    msg=Expected ${count} "${subject}" email(s) to ${address}, found ${found}

Body Of Email
    [Documentation]    The HTML body, falling back to plain text.
    [Arguments]    ${message_id}
    ${response}=    GET On Session    mailpit    /api/v1/message/${message_id}
    ${message}=    Set Variable    ${response.json()}
    ${html}=    Get From Dictionary    ${message}    HTML    default=${EMPTY}
    IF    $html    RETURN    ${html}
    RETURN    ${message}[Text]

Action Token From Email
    [Documentation]    The single-use token in a verify/reset link
    ...    (`{base}/auth/{kind}?[project_id=..&]token=...`).
    [Arguments]    ${address}    ${subject}    ${count}=1
    ${message_id}=    Wait For Email To    ${address}    ${subject}    ${count}
    ${body}=    Body Of Email    ${message_id}
    ${tokens}=    Get Regexp Matches    ${body}    [?&;]token=([A-Za-z0-9._-]+)    1
    Should Not Be Empty    ${tokens}    msg=No token link in the "${subject}" email to ${address}
    RETURN    ${tokens}[0]
