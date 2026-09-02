export const CERTIFICATION_EMISSION_COOLDOWN_MS = 60 * 60 * 1000;

const key = (programId: string) =>
  `univents:certification-emission:${programId}`;

export function certificationEmissionCooldownRemaining(
  programId: string,
  now = Date.now(),
) {
  if (typeof localStorage === "undefined") return 0;
  const expiresAt = Number(localStorage.getItem(key(programId)));
  return Number.isFinite(expiresAt) ? Math.max(0, expiresAt - now) : 0;
}

export function startCertificationEmissionCooldown(programId: string) {
  localStorage.setItem(
    key(programId),
    String(Date.now() + CERTIFICATION_EMISSION_COOLDOWN_MS),
  );
}

export function formatCertificationEmissionCooldown(milliseconds: number) {
  const minutes = Math.ceil(milliseconds / 60_000);
  if (minutes < 60) return `${minutes}min`;
  return `${Math.floor(minutes / 60)}h ${minutes % 60}min`;
}
