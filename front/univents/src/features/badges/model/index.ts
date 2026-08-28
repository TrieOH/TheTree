import z from "zod";

export type {
  BadgeEditionEmission,
  BadgePrintItem,
  BadgeProfileBadge,
  BadgeProfileGroups,
} from "@trieoh/univents-api/schemas";

const bounds = z.object({
  id: z.string(),
  x: z.number(),
  y: z.number(),
  width: z.number().positive(),
  height: z.number().positive(),
});

const richTextRun = z.object({
  text: z.string(),
  bold: z.boolean(),
  italic: z.boolean(),
  underline: z.boolean(),
  color: z.string(),
  fontSize: z.number().min(6).max(400),
  fontFamily: z.string(),
});

const richTextParagraph = z.object({
  align: z.enum(["left", "center", "right", "justify"]),
  lineHeight: z.number().min(0.5).max(4),
  runs: z.array(richTextRun),
});

export const badgeElementSchema = z.discriminatedUnion("type", [
  bounds.extend({
    type: z.literal("text"),
    paragraphs: z.array(richTextParagraph).min(1),
    content: z.string().optional(),
    fontSize: z.number().optional(),
    fontFamily: z.string().optional(),
    color: z.string().optional(),
    fontWeight: z.enum(["normal", "bold"]).optional(),
    align: z.enum(["left", "center", "right"]).optional(),
  }),
  bounds.extend({
    type: z.literal("image"),
    src: z.string(),
    fit: z.enum(["contain", "cover", "fill"]),
    radius: z.number().min(0),
    opacity: z.number().min(0).max(1),
  }),
  bounds.extend({
    type: z.literal("qr"),
    value: z.string(),
    foreground: z.string(),
    background: z.string(),
    style: z.enum(["square", "rounded", "dots"]).default("square"),
  }),
]);

const BADGE_PX_PER_MM = 96 / 25.4;

export const badgeMmToPx = (value: number) => value * BADGE_PX_PER_MM;
export const badgePxToMm = (value: number) =>
  Number((value / BADGE_PX_PER_MM).toFixed(1));
export const MIN_BADGE_CANVAS_SIZE_MM = 26;
export const MIN_BADGE_CANVAS_SIZE_PX = badgeMmToPx(MIN_BADGE_CANVAS_SIZE_MM);

export const badgeDesignSchema = z.object({
  canvas: z.object({
    width: z.number().min(MIN_BADGE_CANVAS_SIZE_PX),
    height: z.number().min(MIN_BADGE_CANVAS_SIZE_PX),
  }),
  backgroundColor: z.string(),
  background: z.string().nullable(),
  elements: z.array(badgeElementSchema),
});

export const badgeTemplateCreateSchema = z.object({
  name: z.string().min(3, "O nome deve ter pelo menos 3 caracteres"),
  ticket_type_id: z.string().nullable(),
  origin: z.enum(["staff"]).nullable().optional(),
  design_data: badgeDesignSchema,
});

export type BadgeElement = z.infer<typeof badgeElementSchema>;
export type BadgeTemplateCreate = z.infer<typeof badgeTemplateCreateSchema>;
export interface BadgeTemplate extends BadgeTemplateCreate {
  id: string;
  edition_id: string;
  created_at: string;
  updated_at: string | null;
}
