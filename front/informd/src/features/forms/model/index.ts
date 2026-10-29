import type {
  CreateFormRequest,
  Form,
  FormStatus,
} from "@trieoh/informd-api/schemas";
import z from "zod";

export type FormStatusI =
  | typeof FormStatus.draft
  | typeof FormStatus.open
  | typeof FormStatus.closed
  | typeof FormStatus.archived;

export const formCreateSchema = z.object({
  title: z
    .string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
}) satisfies z.ZodType<CreateFormRequest>;

export type FormCreateI = CreateFormRequest;

// export type FormI = Form;
export interface FormI extends Omit<Form, "status"> {
  status: FormStatusI;
}
