import { parseUserAgent } from "./parse-user-agent";
import type { DeviceType } from "./device-types";

export interface DeviceInfo {
  device: DeviceType;
  os?: string;
  browser?: string;
}

export function getDeviceInfo(userAgent: string): DeviceInfo {
  const parsed = parseUserAgent(userAgent);

  return {
    device: normalizeDevice(parsed.deviceType),
    os: parsed.os,
    browser: parsed.browser,
  };
}

function normalizeDevice(type?: string): DeviceType {
  if (!type) return "desktop";

  switch (type) {
    case "mobile":
      return "mobile";
    case "tablet":
      return "tablet";
    default:
      return "desktop";
  }
}