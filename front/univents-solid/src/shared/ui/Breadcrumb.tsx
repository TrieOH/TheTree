import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { createMemo } from "solid-js";

import { Breadcrumb as UIBreadcrumb } from "@trieoh/ui-solid";

/** Portuguese labels for the segments we know; anything else keeps its own text. */
const LABELS: Record<string, string> = {
  admin: "Admin",
  events: "Eventos",
  editions: "Edições",
  products: "Produtos",
  purchases: "Compras",
  programs: "Programação",
  members: "Membros",
  signatures: "Assinaturas",
  certifications: "Certificados",
  badges: "Crachás",
  tickets: "Ingressos",
};

const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/**
 * Path-derived crumbs, styled like the React one: uppercase micro-type,
 * chevron separators and id/long segments abbreviated so the strip never grows
 * past the viewport.
 */
export function Breadcrumb(): JSX.Element {
  const location = useLocation();

  const crumbs = createMemo(() => {
    const segments = location().pathname.split("/").filter(Boolean);

    return segments.map((segment, index) => {
      const label = LABELS[segment] ?? segment.charAt(0).toUpperCase() + segment.slice(1);
      const display =
        UUID.test(segment) || label.length > 20
          ? `${label.slice(0, 4)}...${label.slice(-2)}`
          : label;

      return {
        path: `/${segments.slice(0, index + 1).join("/")}`,
        label: display,
      };
    });
  });

  return (
    <UIBreadcrumb
      items={crumbs()}
      renderLink={(item, linkProps) => (
        <Link to={item().path} class={linkProps.class}>{item().label}</Link>
      )}
    />
  );
}
