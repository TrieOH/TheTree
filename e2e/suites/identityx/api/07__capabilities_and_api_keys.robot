*** Settings ***
Documentation       Capabilities name what a credential may do (`resource:action`);
...                 API keys give servers a credential, optionally restricted to some.
Resource            ../../../resources/identityx/steps.resource


*** Test Cases ***
An owner defines a capability
    Given a signed-in user who owns a project
    When they define the capability "intents:checkout"
    Then the capability is created
    And the capability shows in the project's capabilities

A capability is defined only once
    Given a signed-in user who owns a project
    And the project has the capability "intents:checkout"
    When they define the capability "intents:checkout"
    Then the request is rejected as a conflict

Capability names need at least three letters
    Given a signed-in user who owns a project
    When they define the capability "ab:cd"
    Then the request is rejected as invalid
    And the error points at the "resource" field
    And the error points at the "action" field

Project users cannot define capabilities
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they define the capability "intents:refund"
    Then the request is unauthorized

An API key authenticates as the user who made it
    Given a signed-in user who owns a project
    When they create an API key for the project
    Then the raw API key is handed over
    And the API key authenticates as them
    And the API key can list the project's actors

Creating an API key reveals nothing but the raw key
    Given a signed-in user who owns a project
    When they create an API key for the project
    Then the raw API key is handed over
    And only the raw key and public details are returned

A restricted API key carries only its capability
    Given a signed-in user who owns a project
    And the project has the capability "intents:checkout"
    When they create an API key restricted to that capability
    Then the raw API key is handed over
    And the API key carries the capability "intents:checkout"

An expired API key is refused
    Given a signed-in user who owns a project
    When they create an API key that already expired
    Then the raw API key is handed over
    And the API key is refused

A made-up API key is refused
    When someone authenticates with a made-up API key
    Then the request is unauthorized

Project users cannot create API keys
    Given a signed-in user who owns a project
    And a user signed in to that project
    When they create an API key for the project
    Then the request is unauthorized
