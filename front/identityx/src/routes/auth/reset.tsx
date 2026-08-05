import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";

export const Route = createFileRoute("/auth/reset")({
  validateSearch: (search) =>
    z.object({ token: z.string().catch("") }).parse(search),
  beforeLoad: ({ search }) => {
    throw redirect({
      to: "/auth/reset-password",
      search: { token: search.token },
      replace: true,
    });
  },
});
