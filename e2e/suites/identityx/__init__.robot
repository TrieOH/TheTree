*** Settings ***
Documentation       IdentityX end to end: the auth base every other backend
...                 stands on. Runs against the stack `just e2e-up` starts.
Resource            ../../resources/identityx/steps.resource
Suite Setup         The IdentityX stack is ready
