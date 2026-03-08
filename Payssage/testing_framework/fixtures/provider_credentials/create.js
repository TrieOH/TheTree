export const createProviderCredential = (accessToken) => ({
    provider: "mercadopago",
    display_name: "MP Sandbox",
    credentials: {
        access_token: accessToken,
        refresh_token: "",
        provider_user_id: "",
    }
})