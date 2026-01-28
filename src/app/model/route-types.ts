import type { HeaderVariant } from "@/widgets/header/model/header.types";

export interface RouteStaticConfigI {
  header: HeaderVariant;
}

export const RouteComponentTemplate: RouteStaticConfigI = {
  header: "none"
}