const now = new Date();
const inMinutes = (offset) => new Date(now.getTime() + offset * 60 * 1000).toISOString();

export const createEdition = {
    go_auth_event_scope_id: "", // fill in at test time
    type: "year",
    edition_name: "SCTI 2026",
    tagline: "Desenvolvendo o Futuro",
    description: "Te melhorando como dev para o futuro",
    registration_opens_at: inMinutes(5),
    registrations_closes_at: inMinutes(10),
    starts_at: inMinutes(10),
    ends_at: inMinutes(15),
    timezone: "America/Sao_Paulo",
    location_name: "Universidade Estadual Norte Fluminense Darcy Ribeiro",
    location_address: "Alberto Lamego, 2000, Campos dos Goytacazes, RJ",
    logo_url: null,
    banner_url: null,
    contact_email: "sctiuenf@gmail.com",
    contact_phone: "",
    organizer_name: "Giovanna Santos",
};
