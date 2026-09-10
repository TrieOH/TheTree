workspace "TheTree" "C4 Level 1 — System overview of what this repo builds: the four software systems (IdentityX, Univents, Payssage, Informd), their users, and their client relationships to each other and to third-party systems. No implementation details at this level." {

    model {

        // People (product users) ------------------------------------------

        participant = person "Event participant" "Registers for events, buys tickets/products/program spots, checks in, receives badges & certificates, makes doações."
        organizer = person "Event organizer" "Event owner/admin/staff: runs events & editions, sells tickets, takes check-ins, issues badges & certificates."
        respondent = person "Form respondent" "Fills in published Informd forms; may be anonymous."

        // The four software systems this repo builds ------------------------

        identityx = softwareSystem "IdentityX" "Identity & access system: actors, organizations, projects, members & roles, profiles, JWT / API-key auth, OAuth sign-in, key management."
        univents = softwareSystem "Univents" "Event system: editions, ticketing, products & programs, check-in, badges & certificates, purchases."
        payssage = softwareSystem "Payssage" "Payments system: wallets, sellers, payment intents, provider webhooks."
        informd = softwareSystem "Informd" "Forms system: multi-tenant forms, steps/fields, responses."

        // Third-party systems ------------------------------------------------

        googleGithub = softwareSystem "Google / GitHub" "Third-party OAuth identity providers." {
            tags "External"
        }
        mercadoPago = softwareSystem "MercadoPago" "Third-party payment provider: hosted checkout, seller OAuth, payment webhooks." {
            tags "External"
        }

        // people -> systems ----------------------------------------------------

        participant -> univents "registers, buys, checks in, donates"
        participant -> mercadoPago "pays at the hosted checkout"
        organizer -> univents "runs events, editions, check-ins, badges & certificates"
        respondent -> informd "fills in forms"

        // system -> system / system -> third-party ------------------------------

        univents -> identityx "authenticates (client of IdentityX)"
        payssage -> identityx "authenticates (client of IdentityX)"
        informd -> identityx "authenticates (client of IdentityX)"
        univents -> payssage "creates payment intents (client of Payssage)"
        identityx -> googleGithub "delegated sign-in (OAuth2)"
        payssage -> mercadoPago "charges via payment intents"
        mercadoPago -> payssage "delivers payment webhooks"
    }

    views {

        systemLandscape "system-overview" {
            include *
            autoLayout lr
            description "Level 1 — system overview"
        }

        theme default
    }
}
