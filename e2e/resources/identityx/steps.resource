*** Settings ***
Documentation       The English the IdentityX scenarios are written in. Each step
...                 is a Given/When/Then sentence (Robot drops the prefix when it
...                 matches) built from the code-level keywords in api.resource,
...                 ui.resource and mailpit.resource. Suites import only this file.
...
...                 The sentences live in steps/, one file per area:
...                 common (state, the stack, actors, generic outcomes), accounts
...                 (registration, sessions, recovery), tenants (organizations,
...                 projects, members, actors), project_settings (capabilities,
...                 API keys, email templates, terms, profiles, social login) and
...                 app (the browser).
Resource            steps/common.resource
Resource            steps/accounts.resource
Resource            steps/tenants.resource
Resource            steps/project_settings.resource
Resource            steps/app.resource
