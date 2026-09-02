import { describe, expect, it } from "vitest";
import { formatCertificationEmissionCooldown } from "./emission-cooldown";

describe("formatCertificationEmissionCooldown", () => {
  it("rounds up and formats the remaining cooldown", () => {
    expect(formatCertificationEmissionCooldown(59_001)).toBe("1min");
    expect(formatCertificationEmissionCooldown(3_599_001)).toBe("1h 0min");
  });
});
