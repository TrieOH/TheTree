*** Settings ***
Documentation       Sentences about tenants over the API: organizations and their
...                 members, projects and their members, and the actors a project holds.
Library             Collections
Resource            common.resource


*** Keywords ***
# ── organizations ─────────────────────────────────────────────────────────

They create an organization called "${name}"
    ${slug}=    Unique Slug    org
    ${response}=    IDX Create Organization    ${TOKENS}[access_token]    ${name}    ${slug}
    Remember The Response    ${response}

They try to create an organization
    They create an organization called "Not Allowed"

They create another organization with the same slug
    ${response}=    IDX Create Organization    ${TOKENS}[access_token]    E2E Copycat    ${ORGANIZATION}[slug]
    Remember The Response    ${response}

They create an organization with a two-letter name and a one-letter slug
    ${response}=    IDX Create Organization    ${TOKENS}[access_token]    ab    x
    Remember The Response    ${response}

They add the other user to the organization as "${role}"
    ${response}=    IDX Add Organization Member
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${OTHER_EMAIL}    ${role}
    Remember The Response    ${response}

They add the other user to the organization again
    ${response}=    IDX Add Organization Member
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${OTHER_EMAIL}    member
    Remember The Response    ${response}

They add an email nobody registered to the organization
    ${email}=    Unique Email
    ${response}=    IDX Add Organization Member    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${email}
    Remember The Response    ${response}

They remove the other user from the organization
    ${response}=    IDX Remove Organization Member
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${OTHER_EMAIL}
    Remember The Response    ${response}

The other user is a "${role}" of the organization
    [Documentation]    Sets up the membership the scenario starts from.
    ${response}=    IDX Add Organization Member
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${OTHER_EMAIL}    ${role}
    Response Status Should Be    201    ${response}

They list the organization's members
    ${response}=    IDX List Organization Members    ${TOKENS}[access_token]    ${ORGANIZATION}[id]
    Remember The Response    ${response}

They list the organization's projects
    ${response}=    IDX List Organization Projects    ${TOKENS}[access_token]    ${ORGANIZATION}[id]
    Remember The Response    ${response}

They try to add someone to the organization
    [Documentation]    Someone registered, so the answer is about permission, not existence.
    ${email}=    Unique Email
    ${registered}=    IDX Register    ${email}    ${STRONG_PASSWORD}
    Response Status Should Be    201    ${registered}
    ${response}=    IDX Add Organization Member    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${email}
    Remember The Response    ${response}

The member is added
    Response Status Should Be    201    ${RESPONSE}

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
    ${ids}=    Ids Of    ${response}
    List Should Contain Value    ${ids}    ${ORGANIZATION}[id]

The organization shows in the other user's organizations
    ${response}=    IDX List Organizations    ${OTHER_TOKENS}[access_token]
    Response Status Should Be    200    ${response}
    ${ids}=    Ids Of    ${response}
    List Should Contain Value    ${ids}    ${ORGANIZATION}[id]

They are listed as its owner
    Response Status Should Be    200    ${RESPONSE}
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${members}=    Items Of    ${RESPONSE}
    ${roles}=    Evaluate    {m["actor_id"]: m["role"] for m in $members}
    Dictionary Should Contain Item    ${roles}    ${me}[id]    owner

The other user is listed as a "${role}"
    ${response}=    IDX List Organization Members    ${TOKENS}[access_token]    ${ORGANIZATION}[id]
    Response Status Should Be    200    ${response}
    ${other}=    Subject Behind Token    ${OTHER_TOKENS}[access_token]
    ${members}=    Items Of    ${response}
    ${roles}=    Evaluate    {m["actor_id"]: m["role"] for m in $members}
    Dictionary Should Contain Item    ${roles}    ${other}[id]    ${role}

The other user can no longer see the organization's members
    ${response}=    IDX List Organization Members    ${OTHER_TOKENS}[access_token]    ${ORGANIZATION}[id]
    Response Status Should Be    403    ${response}

# ── projects ──────────────────────────────────────────────────────────────

They create a project called "${name}" in their organization
    ${response}=    Create A Project In The Organization    ${name}
    Remember The Response    ${response}

They try to create a project in the organization
    They create a project called "Not Allowed" in their organization

They create a project with the brand slug "${slug}"
    ${response}=    Create A Project In The Organization    E2E Branded    slug=${slug}
    Remember The Response    ${response}

They create another project with the same brand slug
    ${response}=    Create A Project In The Organization    E2E Copycat    slug=${PROJECT}[brand_slug]
    Remember The Response    ${response}

They create another project on the same domain
    ${slug}=    Unique Brand Slug    prj
    ${response}=    IDX Create Organization Project
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    E2E Copycat    ${slug}    domain=${PROJECT}[domain]
    Remember The Response    ${response}

The project belongs to their organization
    Response Status Should Be    201    ${RESPONSE}
    ${project}=    Data Of    ${RESPONSE}
    Should Be Equal    ${project}[organization_id]    ${ORGANIZATION}[id]
    VAR    ${PROJECT}    ${project}    scope=TEST

The project is listed
    Response Status Should Be    200    ${RESPONSE}
    ${ids}=    Ids Of    ${RESPONSE}
    List Should Contain Value    ${ids}    ${PROJECT}[id]

They are listed as the project's owner
    ${response}=    IDX List Organization Project Members
    ...    ${TOKENS}[access_token]    ${ORGANIZATION}[id]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${members}=    Items Of    ${response}
    ${owners}=    Evaluate    [m for m in $members if m["actor_id"] == $me["id"] and m["role"] == "owner"]
    Should Not Be Empty    ${owners}    msg=Owner missing from project members: ${members}
    # Membership starts when the project does, not at Go's zero time.
    ${joined}=    Evaluate    $owners[0]["joined_at"][:4]
    Should Not Be Equal    ${joined}    0001    msg=Owner's joined_at is the zero time: ${owners}[0]

# ── project members ───────────────────────────────────────────────────────

They add the other user to the project as "${role}"
    ${response}=    IDX Add Project Member    ${TOKENS}[access_token]    ${PROJECT}[id]    ${OTHER_EMAIL}    ${role}
    Remember The Response    ${response}

They remove the other user from the project
    ${response}=    IDX Remove Project Member    ${TOKENS}[access_token]    ${PROJECT}[id]    ${OTHER_EMAIL}
    Remember The Response    ${response}

The other user is a member of the project
    ${response}=    IDX List Project Members    ${TOKENS}[access_token]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${other}=    Subject Behind Token    ${OTHER_TOKENS}[access_token]
    ${ids}=    Ids Of    ${response}    actor_id
    List Should Contain Value    ${ids}    ${other}[id]

The other user is not a member of the project
    ${response}=    IDX List Project Members    ${TOKENS}[access_token]    ${PROJECT}[id]
    Response Status Should Be    200    ${response}
    ${other}=    Subject Behind Token    ${OTHER_TOKENS}[access_token]
    ${ids}=    Ids Of    ${response}    actor_id
    List Should Not Contain Value    ${ids}    ${other}[id]

# ── actors ────────────────────────────────────────────────────────────────

The owner lists the project's actors
    ${headers}=    Bearer    ${OWNER_TOKENS}[access_token]
    ${response}=    IDX List Actors    ${headers}    ${PROJECT}[id]
    Remember The Response    ${response}

The owner looks the user up by email
    ${response}=    IDX Get Actor By Email    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${EMAIL}
    Remember The Response    ${response}

The owner looks the user up by id
    ${me}=    Subject Behind Token    ${TOKENS}[access_token]
    ${response}=    IDX Get Actor    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${me}[id]
    Remember The Response    ${response}

The owner looks up an actor that does not exist
    ${id}=    Evaluate    str(uuid.uuid4())    modules=uuid
    ${response}=    IDX Get Actor    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    ${id}
    Remember The Response    ${response}

An outsider lists the project's actors
    Another platform user
    ${headers}=    Bearer    ${OTHER_TOKENS}[access_token]
    ${response}=    IDX List Actors    ${headers}    ${PROJECT}[id]
    Remember The Response    ${response}

The owner provisions a service actor
    ${email}=    Unique Email
    ${response}=    IDX Create Actor    ${OWNER_TOKENS}[access_token]    ${PROJECT}[id]    api_key    service    ${email}
    Remember The Response    ${response}

The user is among the actors
    Response Status Should Be    200    ${RESPONSE}
    ${actors}=    Items Of    ${RESPONSE}
    ${emails}=    Evaluate    [a.get("email") for a in $actors]
    List Should Contain Value    ${emails}    ${EMAIL}

The actor found is the user
    Response Status Should Be    200    ${RESPONSE}
    ${actor}=    Data Of    ${RESPONSE}
    Should Be Equal    ${actor}[email]    ${EMAIL}
    Should Be Equal    ${actor}[project_id]    ${PROJECT}[id]
    Should Be Equal    ${actor}[auth_method]    password
    Should Be Equal    ${actor}[type]    human

A service actor of the project is created
    Response Status Should Be    201    ${RESPONSE}
    ${actor}=    Data Of    ${RESPONSE}
    Should Be Equal    ${actor}[type]    service
    Should Be Equal    ${actor}[auth_method]    api_key
    Should Be Equal    ${actor}[project_id]    ${PROJECT}[id]
