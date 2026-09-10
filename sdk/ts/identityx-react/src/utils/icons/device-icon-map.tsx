import {
  Smartphone,
  Tablet,
  Monitor,
  type LucideIcon
} from "lucide-react";

export const deviceIconMap: Record<string, LucideIcon> = {
  mobile: Smartphone,
  tablet: Tablet,
  desktop: Monitor,
};