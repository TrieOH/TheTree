import { Show, createSignal, onSettled, untrack } from "solid-js";
import type { JSX } from "@solidjs/web";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import CircleUserIcon from "~icons/lucide/circle-user-round";
import GlobeIcon from "~icons/lucide/globe";
import LanguagesIcon from "~icons/lucide/languages";
import LockKeyholeIcon from "~icons/lucide/lock-keyhole";
import Settings2Icon from "~icons/lucide/settings-2";
import InstagramIcon from "~icons/simple-icons/instagram";
import LinkedinIcon from "~icons/simple-icons/linkedin";
import GithubIcon from "~icons/simple-icons/github";
import YoutubeIcon from "~icons/simple-icons/youtube";
import XIcon from "~icons/simple-icons/x";
import TwitterIcon from "~icons/simple-icons/twitter";
import TwitchIcon from "~icons/simple-icons/twitch";
import BlueskyIcon from "~icons/simple-icons/bluesky";
import DiscordIcon from "~icons/simple-icons/discord";
import SaveIcon from "~icons/lucide/save";
import { toast } from "@/shared/ui/toast";
import { Combobox } from "@/shared/ui/Combobox";
import { ProfileImageInput } from "./ProfileImageInput";
import { uploadProfileImage } from "@/features/storage/api";
import {
  asUniventsProfile,
  socialHref,
  type UniventsProfile,
} from "../model/profile-data";

const ArrowLeft =
  ArrowLeftIcon as unknown as () => JSX.Element;
const Save = SaveIcon as unknown as () => JSX.Element;
const CircleUser = CircleUserIcon as unknown as () => JSX.Element;
const Globe = GlobeIcon as unknown as () => JSX.Element;
const Languages = LanguagesIcon as unknown as () => JSX.Element;
const LockKeyhole = LockKeyholeIcon as unknown as () => JSX.Element;
const Settings2 = Settings2Icon as unknown as () => JSX.Element;
const SocialIcons: Record<string, () => JSX.Element> = {
  instagram: InstagramIcon as unknown as () => JSX.Element,
  linkedin: LinkedinIcon as unknown as () => JSX.Element,
  github: GithubIcon as unknown as () => JSX.Element,
  youtube: YoutubeIcon as unknown as () => JSX.Element,
  x: XIcon as unknown as () => JSX.Element,
  twitter: TwitterIcon as unknown as () => JSX.Element,
  twitch: TwitchIcon as unknown as () => JSX.Element,
  bluesky: BlueskyIcon as unknown as () => JSX.Element,
  discord: DiscordIcon as unknown as () => JSX.Element,
};

export function ProfileEditor(props: {
  load: () => Promise<{
    profile: {
      success: boolean;
      data?: { handle?: string | null; pfp_url?: string | null; profile?: Record<string, unknown> };
      message?: string;
    };
  }>;
  save: (
    profile: Record<string, unknown>,
    handle?: string,
  ) => Promise<{ success: boolean; message?: string }>;
  onCancel: () => void;
  onSaved: () => void;
}) {
  const [profile, setProfile] = createSignal<UniventsProfile>({});
  const [handle, setHandle] = createSignal("");
  const [loading, setLoading] = createSignal(true);
  const [saving, setSaving] = createSignal(false);
  const [pendingImages, setPendingImages] = createSignal<Partial<Record<"pfpUrl" | "bannerUrl", File>>>({});
  const [error, setError] = createSignal<string>();

  onSettled(() => {
    void props
      .load()
      .then((result) => {
        if (!result.profile.success || !result.profile.data) {
          setError(
            result.profile.message ?? "Não foi possível carregar o perfil.",
          );
          return;
        }
        setProfile(asUniventsProfile({
          ...result.profile.data.profile,
          ...(result.profile.data.pfp_url !== undefined && {
            pfpUrl: result.profile.data.pfp_url,
          }),
        }));
        setHandle(result.profile.data.handle ?? "");
      })
      .catch((cause) =>
        setError(
          cause instanceof Error
            ? cause.message
            : "Não foi possível carregar o perfil.",
        ),
      )
      .finally(() => setLoading(false));
  });

  const update = (field: string, value: unknown) =>
    setProfile((current) => ({ ...current, [field]: value }));
  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setSaving(true);
    try {
      const nextProfile = untrack(() => ({ ...profile() }));
      const uploads = untrack(() => Object.entries(pendingImages())) as ["pfpUrl" | "bannerUrl", File][];
      const results = await Promise.allSettled(uploads.map(async ([field, file]) => [field, await uploadProfileImage(file, field)] as const));
      const failed: string[] = [];
      results.forEach((result, index) => {
        if (result.status === "fulfilled") nextProfile[result.value[0]] = result.value[1];
        else failed.push(uploads[index]?.[0] === "pfpUrl" ? "foto" : "banner");
      });
      const result = await props.save(
        nextProfile,
        handle().trim().replace(/^@/, ""),
      );
      if (!result.success)
        throw new Error(result.message ?? "Não foi possível salvar o perfil.");
      toast.success("Perfil atualizado");
      if (failed.length) toast.warning(`Perfil salvo, mas falhou o upload de: ${failed.join(" e ")}.`);
      props.onSaved();
    } catch (cause) {
      toast.error(
        cause instanceof Error
          ? cause.message
          : "Não foi possível salvar o perfil.",
      );
    } finally {
      setSaving(false);
    }
  };

  return (
    <Show
      when={!loading()}
      fallback={<div class="min-h-dvh animate-pulse bg-muted" />}
    >
      <Show
        when={!error()}
        fallback={
          <main class="mx-auto max-w-3xl px-4 py-16 text-center text-muted-foreground">
            {error()}
          </main>
        }
      >
        <form onSubmit={submit} class="min-h-dvh bg-background pb-28">
          <section class="border-b border-border bg-card shadow-md">
            <div class="relative h-40 overflow-hidden bg-linear-to-br from-primary/40 via-primary/15 to-muted sm:h-48 md:h-56">
              <ProfileImageInput
                label="Banner"
                currentUrl={profile().bannerUrl}
                variant="banner"
                onSelect={(file) => setPendingImages((value) => ({ ...value, bannerUrl: file }))}
              />
              <button
                type="button"
                onClick={() => props.onCancel()}
                class="absolute left-4 top-4 inline-flex size-10 items-center justify-center rounded-full bg-background shadow-lg"
                aria-label="Voltar ao perfil"
              >
                <ArrowLeft />
              </button>
            </div>
            <div class="mx-auto max-w-7xl px-4">
              <div class="relative pb-5 md:pb-6">
                <div class="-mt-12 md:-mt-16">
                  <ProfileImageInput
                    label="Foto do perfil"
                    currentUrl={profile().pfpUrl}
                    variant="avatar"
                    onSelect={(file) => setPendingImages((value) => ({ ...value, pfpUrl: file }))}
                  />
                </div>
                <div class="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                  <Field
                    label="Nome de usuário"
                    value={handle()}
                    onInput={setHandle}
                    placeholder="Escolha um nome de usuário único"
                    autocomplete="username"
                  />
                  <Field
                    label="Nome Social"
                    value={profile().preferredName ?? ""}
                    onInput={(value) => update("preferredName", value)}
                    placeholder="Como você prefere ser chamado?"
                    autocomplete="nickname"
                  />
                  <Field
                    label="Nome Civil"
                    value={profile().legalName ?? ""}
                    onInput={(value) => update("legalName", value)}
                    autocomplete="name"
                    placeholder="Digite seu nome civil completo"
                  />
                  <Field
                    label="Função"
                    value={profile().role ?? ""}
                    onInput={(value) => update("role", value)}
                    placeholder="Ex.: Organizador de eventos"
                    autocomplete="organization-title"
                  />
                  <Field
                    label="Organização"
                    value={profile().organization ?? ""}
                    onInput={(value) => update("organization", value)}
                    placeholder="Ex.: Univents"
                    autocomplete="organization"
                  />
                  <PronounsField
                    label="Pronomes"
                    value={profile().pronouns ?? ""}
                    onInput={(value) => update("pronouns", value)}
                  />
                </div>
              </div>
            </div>
          </section>
          <div class="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:mt-5 md:grid-cols-[minmax(0,1fr)_360px] md:gap-5">
            <div class="space-y-5">
              <Card title="Sobre mim" icon={<CircleUser />}>
                <textarea
                  value={profile().aboutMe ?? ""}
                  onInput={(event) =>
                    update("aboutMe", event.currentTarget.value)
                  }
                  placeholder="Conte um pouco sobre você, sua experiência e seus interesses…"
                  class="min-h-32 w-full resize-y rounded-md border border-border bg-background p-3 text-sm outline-none focus:border-primary"
                />
              </Card>
              <Card title="Idiomas" icon={<Languages />}>
                <LanguagesField
                  value={(profile().languages ?? []).join(", ")}
                  onChange={(value) => update("languages", value)}
                />
              </Card>
              <Card title="Contato" icon={<Globe />}>
                <div class="grid gap-4 sm:grid-cols-2">
                  <Field
                    label="Website"
                    value={profile().website ?? ""}
                    onInput={(value) => update("website", value)}
                    onBlur={(value) => {
                      const trimmed = value.trim();
                      if (trimmed && !/^https?:\/\//i.test(trimmed)) {
                        update("website", `https://${trimmed}`);
                      }
                    }}
                    placeholder="https://seusite.com"
                    type="url"
                  />
                  <Field
                    label="E-mail"
                    value={profile().contactEmail ?? ""}
                    onInput={(value) => update("contactEmail", value)}
                    placeholder="voce@exemplo.com"
                    type="email"
                    autocomplete="email"
                  />
                </div>
                <SocialFields
                  value={profile().socials ?? {}}
                  onChange={(value) => update("socials", value)}
                />
              </Card>
            </div>
            <aside class="space-y-5">
              <Card title="Privacidade" icon={<LockKeyhole />}>
                <Toggle
                  label="Ocultar nome legal"
                  checked={profile().visibility?.hideLegalName === true}
                  onChange={(checked) =>
                    update("visibility", {
                      ...profile().visibility,
                      hideLegalName: checked,
                    })
                  }
                />
              </Card>
              <Card title="Outros detalhes" icon={<Settings2 />}>
                <TimezoneField
                  value={String(profile().timezone ?? "")}
                  onInput={(value) => update("timezone", value)}
                />
              </Card>
            </aside>
          </div>
          <div class="fixed inset-x-0 bottom-0 z-50 border-t border-border bg-background/95 p-3 shadow-[0_-8px_30px_rgb(0_0_0/0.08)] backdrop-blur">
            <div class="mx-auto flex max-w-7xl justify-between gap-3 px-1 md:px-4">
              <button
                type="button"
                onClick={() => props.onCancel()}
                class="inline-flex h-10 items-center gap-2 rounded-md px-4 text-sm hover:bg-muted"
              >
                <ArrowLeft />
                Cancelar
              </button>
              <button
                type="submit"
                disabled={saving()}
                class="inline-flex h-10 items-center gap-2 rounded-md bg-primary px-5 text-sm text-primary-foreground disabled:opacity-50"
              >
                <Save />
                {saving() ? "Salvando..." : "Salvar alterações"}
              </button>
            </div>
          </div>
        </form>
      </Show>
    </Show>
  );
}

function Card(props: { title: string; icon?: JSX.Element; children: unknown }) {
  return (
    <section class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
      <h2 class="mb-4 flex items-center gap-2 text-base font-semibold [&>svg]:size-4 [&>svg]:text-muted-foreground">
        {props.icon}
        {props.title}
      </h2>
      {props.children}
    </section>
  );
}
function Field(props: {
  label: string;
  value: string;
  onInput: (value: string) => void;
  onBlur?: (value: string) => void;
  placeholder?: string;
  autocomplete?: string;
  type?: "text" | "url" | "email";
}) {
  return (
    <label class="block space-y-2.5 text-sm font-medium">
      {props.label}
      <input
        type={props.type ?? "text"}
        value={props.value}
        placeholder={props.placeholder}
        autocomplete={props.autocomplete}
        onInput={(event) => props.onInput(event.currentTarget.value)}
        onBlur={(event) => props.onBlur?.(event.currentTarget.value)}
        class="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-primary"
      />
    </label>
  );
}

function PronounsField(props: {
  label: string;
  value: string;
  onInput: (value: string) => void;
}) {
  const options = [
    "Ele/dele",
    "Ela/dela",
    "Elu/delu",
    "Eles/deles",
    "Ela/ele",
    "Prefiro não informar",
    "Outros",
  ].map((value) => ({ value, label: value }));
  const initialValue = untrack(() => props.value);
  const [custom, setCustom] = createSignal(
    options.some((option) => option.value === initialValue) ? "" : initialValue,
  );
  const [otherSelected, setOtherSelected] = createSignal(
    initialValue === "Outros" ||
    !options.some((option) => option.value === initialValue),
  );
  return (
    <label class="block space-y-2.5 text-sm font-medium">
      {props.label}
      <Combobox
        value={otherSelected() ? "Outros" : props.value}
        options={options}
        placeholder="Selecione seus pronomes"
        onChange={(value) => {
          if (value === "Outros") {
            setOtherSelected(true);
            props.onInput(custom() || "Outros");
          } else {
            setOtherSelected(false);
            props.onInput(value);
          }
        }}
      />
      <Show when={otherSelected()}>
        <input
          value={custom()}
          placeholder="Digite seus pronomes"
          onInput={(event) => {
            setCustom(event.currentTarget.value);
            props.onInput(event.currentTarget.value || "Outros");
          }}
          class="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-primary"
        />
      </Show>
    </label>
  );
}

function LanguagesField(props: {
  value: string;
  onChange: (value: string[]) => void;
}) {
  const [draft, setDraft] = createSignal("");
  const items = () =>
    props.value
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean);
  const add = (value: string) => {
    const item = value.trim();
    if (!item) return;
    props.onChange([...items(), item]);
    setDraft("");
  };
  const input = (value: string) => {
    const parts = value.split(/[,;\n]+/);
    if (parts.length > 1) {
      parts.slice(0, -1).forEach(add);
      setDraft(parts.at(-1) ?? "");
    } else setDraft(value);
  };
  return (
    <div class="space-y-2.5">
      <Show when={items().length > 0}>
        <div class="flex flex-wrap gap-1.5">
          {items().map((item, index) => (
            <button
              type="button"
              class="rounded-full border border-border bg-muted px-2.5 py-1 text-xs"
              onClick={() =>
                props.onChange(
                  items().filter((_, entryIndex) => entryIndex !== index),
                )
              }
              title="Remover idioma"
            >
              {item} ×
            </button>
          ))}
        </div>
      </Show>
      <input
        value={draft()}
        placeholder="Digite um idioma e pressione Enter"
        autocomplete="off"
        onInput={(event) => input(event.currentTarget.value)}
        onBlur={() => add(draft())}
        onKeyDown={(event) => {
          if (event.key === "Enter") {
            event.preventDefault();
            add(draft());
          }
        }}
        class="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none placeholder:text-muted-foreground focus:border-primary"
      />
    </div>
  );
}

function SocialFields(props: {
  value: Record<string, string | null | undefined>;
  onChange: (value: Record<string, string>) => void;
}) {
  const networks = [
    ["instagram", "Instagram", "Ex.: @univents ou instagram.com/univents"],
    ["linkedin", "LinkedIn", "Ex.: trieoh ou linkedin.com/in/trieoh"],
    ["github", "GitHub", "Ex.: trieoh ou github.com/trieoh"],
    ["youtube", "YouTube", "Ex.: @univents ou youtube.com/@univents"],
    ["x", "X", "Ex.: @univents ou x.com/univents"],
    ["twitter", "Twitter", "Ex.: @univents ou twitter.com/univents"],
    ["twitch", "Twitch", "Ex.: univents ou twitch.tv/univents"],
    ["bluesky", "Bluesky", "Ex.: univents.bsky.social"],
    ["discord", "Discord", "Ex.: seu ID ou discord.com/users/seu-id"],
  ];
  const updateNetwork = (network: string, value: string) => {
    props.onChange({
      ...props.value,
      [network]: value,
    } as Record<string, string>);
  };
  const normalize = (value: string) =>
    value
      .trim()
      .replace(/^https?:\/\//i, "")
      .replace(/^@/, "")
      .replace(/\/$/, "");
  return (
    <div class="space-y-3">
      <div class="mt-5 grid gap-3 sm:grid-cols-2">
        {networks.map(([key, label, placeholder]) => (
          <label class="min-w-0 space-y-1.5">
            <div class="relative">
              <span class="pointer-events-none absolute inset-y-0 left-0 flex w-10 items-center justify-center rounded-l-md border border-border bg-primary text-xs font-bold text-primary-foreground [&>svg]:size-4">
                {renderSocialIcon(key, label)}
              </span>
              <input
                value={props.value[key] ?? ""}
                placeholder={placeholder}
                autocomplete="off"
                aria-label={label}
                onInput={(event) =>
                  updateNetwork(key, event.currentTarget.value)
                }
                onBlur={(event) =>
                  updateNetwork(key, normalize(event.currentTarget.value))
                }
                class="h-10 w-full min-w-0 rounded-md border border-border bg-background pl-12 pr-3 text-sm outline-none focus:border-primary"
              />
            </div>
            {props.value[key] && (
              <p
                class="truncate text-right text-xs text-muted-foreground"
                title={socialHref(key, normalize(props.value[key] ?? ""))}
              >
                {socialHref(key, normalize(props.value[key] ?? ""))}
              </p>
            )}
          </label>
        ))}
      </div>
    </div>
  );
}

function renderSocialIcon(key: string, label: string) {
  const Icon = SocialIcons[key];
  return Icon ? <Icon /> : label.slice(0, 2).toUpperCase();
}

function Toggle(props: {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}) {
  return (
    <label class="flex items-center justify-between gap-4 rounded-md border border-border bg-background p-3 text-sm">
      <span>{props.label}</span>
      <input
        type="checkbox"
        checked={props.checked}
        onChange={(event) => props.onChange(event.currentTarget.checked)}
        class="size-4 accent-primary"
      />
    </label>
  );
}

function TimezoneField(props: {
  value: string;
  onInput: (value: string) => void;
}) {
  const zones =
    typeof Intl.supportedValuesOf === "function"
      ? Intl.supportedValuesOf("timeZone")
      : ["America/Sao_Paulo", "UTC"];
  const options = zones.map((zone) => ({
    value: zone,
    label: zone.replaceAll("_", " "),
  }));
  return (
    <label class="block space-y-1.5 text-sm font-medium">
      Fuso horário
      <Combobox
        value={props.value}
        options={options}
        placeholder="Selecione o fuso horário"
        searchPlaceholder="Buscar cidade ou fuso…"
        onChange={props.onInput}
      />
    </label>
  );
}
