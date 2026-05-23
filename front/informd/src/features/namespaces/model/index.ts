import type { CreateNamespaceRequest, Namespace } from "@trieoh/informd-models";
import z from "zod";

export const namespaceCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
}) satisfies z.ZodType<CreateNamespaceRequest>;

export type NamespaceCreateI = CreateNamespaceRequest;

export type NamespaceI = Namespace;