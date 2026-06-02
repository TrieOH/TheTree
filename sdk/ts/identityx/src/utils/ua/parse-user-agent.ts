import { UAParser } from "ua-parser-js";

export function parseUserAgent(userAgent: string) {
  const { device, os, browser } =  UAParser(userAgent);

  return {
    deviceType: device.type ?? "desktop",
    os: os.name,
    browser: browser.name,
  };
}
