import {
  Schema,
  createSchema,
  positive,
  nonNegative,
  min,
  max,
  minLength,
  type Infer,
} from "@/shared/lib/effect-schema";

export type {
  BadgeEditionEmission,
  BadgePrintItem,
  BadgeProfileBadge,
  BadgeProfileGroups,
} from "@trieoh/univents-api/schemas";

const boundsFields = {
  id: Schema.String,
  x: Schema.Number,
  y: Schema.Number,
  width: Schema.Number.pipe(positive()),
  height: Schema.Number.pipe(positive()),
};

const richTextRun = Schema.Struct({
  text: Schema.String,
  bold: Schema.Boolean,
  italic: Schema.Boolean,
  underline: Schema.Boolean,
  color: Schema.String,
  fontSize: Schema.Number.pipe(min(6), max(400)),
  fontFamily: Schema.String,
});

const richTextParagraph = Schema.Struct({
  align: Schema.Literals(["left", "center", "right", "justify"]),
  lineHeight: Schema.Number.pipe(min(0.5), max(4)),
  runs: Schema.Array(richTextRun),
});

const textBadgeElement = Schema.Struct({
  ...boundsFields,
  type: Schema.Literal("text"),
  paragraphs: Schema.Array(richTextParagraph),
  content: Schema.optional(Schema.String),
  fontSize: Schema.optional(Schema.Number),
  fontFamily: Schema.optional(Schema.String),
  color: Schema.optional(Schema.String),
  fontWeight: Schema.optional(Schema.Literals(["normal", "bold"])),
  align: Schema.optional(Schema.Literals(["left", "center", "right"])),
});

const imageBadgeElement = Schema.Struct({
  ...boundsFields,
  type: Schema.Literal("image"),
  src: Schema.String,
  fit: Schema.Literals(["contain", "cover", "fill"]),
  radius: Schema.Number.pipe(nonNegative()),
  opacity: Schema.Number.pipe(min(0), max(1)),
});

const qrBadgeElement = Schema.Struct({
  ...boundsFields,
  type: Schema.Literal("qr"),
  value: Schema.String,
  foreground: Schema.String,
  background: Schema.String,
  style: Schema.optional(Schema.Literals(["square", "rounded", "dots"])),
});

export const badgeElementSchema = createSchema(
  Schema.Union([textBadgeElement, imageBadgeElement, qrBadgeElement]),
);

const BADGE_PX_PER_MM = 96 / 25.4;

export const badgeMmToPx = (value: number) => value * BADGE_PX_PER_MM;
export const badgePxToMm = (value: number) =>
  Number((value / BADGE_PX_PER_MM).toFixed(1));
export const MIN_BADGE_CANVAS_SIZE_MM = 26;
export const MIN_BADGE_CANVAS_SIZE_PX = badgeMmToPx(MIN_BADGE_CANVAS_SIZE_MM);

export const badgeDesignSchema = createSchema(
  Schema.Struct({
    canvas: Schema.Struct({
      width: Schema.Number.pipe(min(MIN_BADGE_CANVAS_SIZE_PX)),
      height: Schema.Number.pipe(min(MIN_BADGE_CANVAS_SIZE_PX)),
    }),
    backgroundColor: Schema.String,
    background: Schema.NullOr(Schema.String),
    elements: Schema.Array(badgeElementSchema),
  }),
);

export const badgeTemplateCreateSchema = createSchema(
  Schema.Struct({
    name: Schema.String.pipe(minLength(3)),
    ticket_type_id: Schema.NullOr(Schema.String),
    origin: Schema.optional(Schema.NullOr(Schema.Literal("staff"))),
    design_data: badgeDesignSchema,
  }),
);

export type BadgeElement = Infer<typeof badgeElementSchema>;
export type BadgeDesign = Infer<typeof badgeDesignSchema>;
export type BadgeTemplateCreate = Infer<typeof badgeTemplateCreateSchema>;
export type BadgeTemplate = BadgeTemplateCreate & {
  id: string;
  edition_id: string;
  created_at: string;
  updated_at: string | null;
};

export * from "./profile-badges";
