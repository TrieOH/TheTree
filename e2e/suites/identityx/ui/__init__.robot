*** Settings ***
Documentation       IdentityX through its front (front/identityx) in Chromium.
Resource            ../../../resources/identityx/steps.resource
Test Tags           ui
Suite Setup         Open The Browser
Suite Teardown      Close Browser    ALL
Test Setup          Open A Fresh Browser Context
Test Teardown       Close Context
