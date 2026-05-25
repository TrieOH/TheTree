import type {
  CreateFormRequest,
  Form,
  FormStatusArchived,
  FormStatusClosed,
  FormStatusDraft,
  FormStatusOpen
} from "@trieoh/informd-models";
import z from "zod";

export type FormStatusI =
  | typeof FormStatusDraft
  | typeof FormStatusOpen
  | typeof FormStatusClosed
  | typeof FormStatusArchived;

export const formCreateSchema = z.object({
  title: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
}) satisfies z.ZodType<CreateFormRequest>;

export type FormCreateI = CreateFormRequest;

// export type FormI = Form;
export interface FormI
  extends Omit<Form, "status"> {
  status: FormStatusI;
}