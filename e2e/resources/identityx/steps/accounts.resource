*** Settings ***
Documentation       Sentences about accounts over the API: the service itself,
...                 registration, email verification, sessions and password resets.
Library             Collections
Library             String
Resource            common.resource


*** Keywords ***
# ── the service ───────────────────────────────────────────────────────────

They check the service health
    ${response}=    IDX Get Health
    Remember The Response    ${response}

Someone asks whether setup is still pending
    ${response}=    IDX Get Setup Status
    Remember The Response    ${response}

Someone tries to set up the platform again
    ${email}=    Unique Email
    ${response}=    IDX Set Up Platform    ${email}    ${STRONG_PASSWORD}
    Remember The Response    ${response}

Someone fetches the API specification
    ${response}=    IDX Get OpenAPI Spec
    Remember The Response    ${response}

Someone fetches the platform signing keys
    ${response}=    IDX Get JWKS
    Remember The Response    ${response}

Someone fetches that project's signing keys
    ${response}=    IDX Get JWKS    project_id=${PROJECT}[id]
    Remember The Response    ${response}

The service reports it is healthy
    Response Status Should Be    200    ${RESPONSE}
    ${body}=    Data Of    ${RESPONSE}
    Should Be Equal    ${body}[status]    ok
    Should Be Equal    ${body}[service]    IdentityX

Setup is reported as complete
    Response Status Should Be    409    ${RESPONSE}

The specification describes the IdentityX API
    Response Status Should Be    200    ${RESPONSE}
    Should Start With    ${RESPONSE.text}    openapi: 3.1
    Should Contain    ${RESPONSE.text}    title: IdentityX API

They get Ed25519 public keys
    [Documentation]    A raw RFC 7517 document, not the fun envelope.
    Response Status Should Be    200    ${RESPONSE}
    ${keys}=    Set Variable    ${RESPONSE.json()}[keys]
    Should Not Be Empty    ${keys}
    FOR    ${key}    IN    @{keys}
        Should Be Equal    ${key}[kty]    OKP
        Should Be Equal    ${key}[crv]    Ed25519
        Should Not Be Empty    ${key}[kid]
        Dictionary Should Not Contain Key    ${key}    d    msg=A private key part leaked into the JWKS
    END

The keys differ from the platform's
    ${platform}=    IDX Get JWKS
    ${platform_kids}=    Evaluate    [k["kid"] for k in $platform.json()["keys"]]
    ${project_kids}=    Evaluate    [k["kid"] for k in $RESPONSE.json()["keys"]]
    FOR    ${kid}    IN    @{project_kids}
        Should Not Contain    ${platform_kids}    ${kid}
        ...    msg=The platform keyset also publishes project key ${kid}    values=${False}
    END

Their access token was signed with one of those keys
    ${header}=    Evaluate
    ...    json.loads(base64.urlsafe_b64decode($TOKENS["access_token"].split(".")[0] + "=="))
    ...    modules=json,base64
    Should Be Equal    ${header}[alg]    EdDSA
    ${kids}=    Evaluate    [k["kid"] for k in $RESPONSE.json()["keys"]]
    List Should Contain Value    ${kids}    ${header}[kid]

# ── registration and verification ─────────────────────────────────────────

They register on the platform
    ${response}=    IDX Register    ${EMAIL}    ${PASSWORD}
    Remember The Response    ${response}

They register on the platform again
    They register on the platform

They register on the platform as "${email}"
    ${response}=    IDX Register    ${email}    ${PASSWORD}
    Remember The Response    ${response}

They confirm their email with the link from the inbox
    ${token}=    Action Token From Email    ${EMAIL}    ${VERIFY_SUBJECT}
    ${response}=    IDX Verify Email    ${token}
    Remember The Response    ${response}
    VAR    ${USED_VERIFY_TOKEN}    ${token}    scope=TEST

They open the same verification link again
    ${response}=    IDX Verify Email    ${USED_VERIFY_TOKEN}
    Remember The Response    ${response}

Someone verifies an email with a forged token
    ${response}=    IDX Verify Email    not-a-real-token
    Remember The Response    ${response}

They ask for a new verification email
    ${response}=    IDX Resend Verification    ${EMAIL}
    Remember The Response    ${response}

Someone asks for a verification email for an address nobody registered
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    ${response}=    IDX Resend Verification    ${EMAIL}
    Remember The Response    ${response}

A verification email arrives in their inbox
    Wait For Email To    ${EMAIL}    ${VERIFY_SUBJECT}

A second verification email arrives in their inbox
    Wait For Email To    ${EMAIL}    ${VERIFY_SUBJECT}    count=2

No verification email is sent to them
    No Email Should Arrive To    ${EMAIL}    ${VERIFY_SUBJECT}

No second verification email is sent to them
    Sleep    ${EMAIL_GRACE}
    ${messages}=    Emails To    ${EMAIL}    ${VERIFY_SUBJECT}
    Length Should Be    ${messages}    1    msg=An already-verified account got another verification email

Their identity shows the email as verified
    ${verified_at}=    Verified At    ${TOKENS}[access_token]
    Should Not Be Equal    ${verified_at}    ${None}    msg=Email is still unverified

Their identity shows the email as unverified
    ${verified_at}=    Verified At    ${TOKENS}[access_token]
    Should Be Equal    ${verified_at}    ${None}    msg=Email is already verified

# ── sessions ──────────────────────────────────────────────────────────────

They sign in
    ${response}=    IDX Login    ${EMAIL}    ${PASSWORD}
    Remember The Response    ${response}

They sign in with a wrong password
    ${response}=    IDX Login    ${EMAIL}    Wr0ng-${PASSWORD}
    Remember The Response    ${response}

They sign in with their email in capitals
    ${upper}=    Convert To Upper Case    ${EMAIL}
    ${response}=    IDX Login    ${upper}    ${PASSWORD}
    Remember The Response    ${response}

They sign in to the platform instead of their project
    ${response}=    IDX Login    ${EMAIL}    ${PASSWORD}
    Remember The Response    ${response}

They sign in to that project with their platform account
    ${response}=    IDX Login    ${EMAIL}    ${PASSWORD}    project_id=${PROJECT}[id]
    Remember The Response    ${response}

Someone signs in with an email nobody registered
    ${email}=    Unique Email
    ${response}=    IDX Login    ${email}    ${STRONG_PASSWORD}
    Remember The Response    ${response}

They refresh their session
    VAR    ${OLD_TOKENS}    ${TOKENS}    scope=TEST
    ${response}=    IDX Refresh    ${TOKENS}[refresh_token]
    Remember The Response    ${response}

They try to refresh with their access token
    [Documentation]    An access token is not a refresh token: the purposes must not mix.
    ${response}=    IDX Refresh    ${TOKENS}[access_token]
    Remember The Response    ${response}

Someone refreshes with a malformed token
    ${response}=    IDX Refresh    not-a-jwt
    Remember The Response    ${response}

They log out
    ${response}=    IDX Logout    ${TOKENS}[access_token]    ${TOKENS}[refresh_token]
    Remember The Response    ${response}

Someone asks who they are without credentials
    ${response}=    IDX Introspect Anonymously
    Remember The Response    ${response}

Someone asks who they are with a forged token
    ${response}=    IDX Introspect With Token    ${TOKENS}[access_token]x
    Remember The Response    ${response}

They receive a fresh token pair
    Response Status Should Be    200    ${RESPONSE}
    ${tokens}=    Data Of    ${RESPONSE}
    Should Not Be Empty    ${tokens}[access_token]
    Should Not Be Empty    ${tokens}[refresh_token]
    VAR    ${TOKENS}    ${tokens}    scope=TEST

The access token expires before the refresh token
    ${access}=    Evaluate    datetime.datetime.fromisoformat($TOKENS["access_expires_at"])    modules=datetime
    ${refresh}=    Evaluate    datetime.datetime.fromisoformat($TOKENS["refresh_expires_at"])    modules=datetime
    Should Be True    $access < $refresh
    ...    msg=Access token (${access}) outlives the refresh token (${refresh})

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

Their identity belongs to no project
    ${subject}=    Subject Behind Token    ${TOKENS}[access_token]
    ${project_id}=    Get From Dictionary    ${subject}    project_id    default=${None}
    Should Be Equal    ${project_id}    ${None}

Their identity belongs to that project
    ${subject}=    Subject Behind Token    ${TOKENS}[access_token]
    Should Be Equal    ${subject}[project_id]    ${PROJECT}[id]

They are identified as a human using a token
    ${response}=    IDX Introspect With Token    ${TOKENS}[access_token]
    Response Status Should Be    200    ${response}
    ${identity}=    Data Of    ${response}
    Should Be Equal    ${identity}[subject][type]    human
    Should Be Equal    ${identity}[cred][type]    token

# ── password reset ────────────────────────────────────────────────────────

They ask for a password reset
    ${response}=    IDX Forgot Password    ${EMAIL}
    Remember The Response    ${response}

Someone asks for a password reset for an address nobody registered
    ${email}=    Unique Email
    VAR    ${EMAIL}    ${email}    scope=TEST
    ${response}=    IDX Forgot Password    ${EMAIL}
    Remember The Response    ${response}

They choose the new password "${new_password}" with the link from the inbox
    ${token}=    Action Token From Email    ${EMAIL}    ${RESET_SUBJECT}
    ${response}=    IDX Reset Password    ${token}    ${new_password}
    Remember The Response    ${response}
    VAR    ${USED_RESET_TOKEN}    ${token}    scope=TEST

The same reset link is used again for "${new_password}"
    ${response}=    IDX Reset Password    ${USED_RESET_TOKEN}    ${new_password}
    Remember The Response    ${response}

Someone resets a password with a forged token
    ${response}=    IDX Reset Password    not-a-real-token    ${STRONG_PASSWORD}
    Remember The Response    ${response}

They use the verification link to reset their password
    [Documentation]    Tokens are purpose-scoped: a verify token must not reset a password.
    ${token}=    Action Token From Email    ${EMAIL}    ${VERIFY_SUBJECT}
    ${response}=    IDX Reset Password    ${token}    N3w-Passw0rd!
    Remember The Response    ${response}

A password reset email arrives in their inbox
    Wait For Email To    ${EMAIL}    ${RESET_SUBJECT}

No password reset email is sent to them
    No Email Should Arrive To    ${EMAIL}    ${RESET_SUBJECT}

They can sign in with the password "${password}"
    ${response}=    IDX Login    ${EMAIL}    ${password}
    Response Status Should Be    200    ${response}

They can no longer sign in with the password "${password}"
    ${response}=    IDX Login    ${EMAIL}    ${password}
    Response Status Should Be    401    ${response}
