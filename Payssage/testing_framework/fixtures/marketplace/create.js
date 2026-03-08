export const createMarketplaceConfig = (credentialId, feeBps = 500) => ({
    credential_id: credentialId,
    fee_bps: feeBps,
})