import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";

export const Route = createFileRoute("/auth/verify")({
  validateSearch: (search) =>
    z.object({ token: z.string().catch("") }).parse(search),
  beforeLoad: ({ search }) => {
    throw redirect({
      to: "/auth/verify-email",
      search: { token: search.token },
      replace: true,
    });
  },
});
