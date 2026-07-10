import { useState } from "react";
import type { FieldFormApi } from "@/widgets/multi-step-form/model/types";
import type { EventCreateInputI, SocialPlatform } from "../../model";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";

const socialPlatformMeta: Record<SocialPlatform, { label: string; placeholder: string }> = {
  website: { label: "Website", placeholder: "https://seu-evento.com" },
  instagram: { label: "Instagram", placeholder: "https://instagram.com/seu-evento" },
  linkedin: { label: "LinkedIn", placeholder: "https://linkedin.com/company/seu-evento" },
  twitter: { label: "Twitter", placeholder: "https://twitter.com/seu-evento" },
};

/**
 * Bespoke composite: the four toggle buttons + a single URL input that
 * targets whichever platform is selected. Doesn't fit a plain "text"
 * field, so it's registered as a "custom" field instead of forcing the
 * generic engine to understand toggle groups.
 */
export function SocialLinksField({ form }: { form: FieldFormApi<EventCreateInputI> }) {
  const [selected, setSelected] = useState<SocialPlatform>("instagram");
  const meta = socialPlatformMeta[selected];
  const fieldName = `social_links.${selected}` as const;
  const error = form.formState.errors.social_links?.[selected]?.message;

  return (
    <div className="space-y-4">
      <div>
        <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
          Redes sociais
        </span>
        <div className="mt-2 flex flex-wrap gap-2">
          {(Object.keys(socialPlatformMeta) as SocialPlatform[]).map((platform) => {
            const isActive = platform === selected;
            return (
              <Button
                key={platform}
                type="button"
                variant={isActive ? "default" : "outline"}
                size="sm"
                onClick={() => setSelected(platform)}
              >
                {socialPlatformMeta[platform].label}
              </Button>
            );
          })}
        </div>
      </div>

      <div className="space-y-1.5">
        <Label
          htmlFor={fieldName}
          className="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
        >
          {meta.label} URL
        </Label>
        <Input
          id={fieldName}
          placeholder={meta.placeholder}
          {...form.register(fieldName)}
        />
        {error ? <p className="text-xs text-destructive">{error}</p> : null}
      </div>
    </div>
  );
}
