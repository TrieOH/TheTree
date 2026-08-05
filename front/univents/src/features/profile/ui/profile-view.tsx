import { Link } from "@tanstack/react-router";
import type { ActorProfile } from "@trieoh/identityx-sdk-ts";
import { Calendar, CircleAlert, Copy, Globe, Mail } from "lucide-react";
import { useEffect, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import blueskyIcon from "@/shared/ui/social-icons/assets/bluesky.svg";
import discordIcon from "@/shared/ui/social-icons/assets/discord.svg";
import githubIcon from "@/shared/ui/social-icons/assets/github.svg";
import instagramIcon from "@/shared/ui/social-icons/assets/instagram.svg";
import linkedinIcon from "@/shared/ui/social-icons/assets/linkedin.svg";
import twitchIcon from "@/shared/ui/social-icons/assets/twitch.svg";
import xIcon from "@/shared/ui/social-icons/assets/x.svg";
import youtubeIcon from "@/shared/ui/social-icons/assets/youtube.svg";
import {
  asUniventsProfile,
  profileCompleteness,
  profileDisplayName,
  socialHref,
} from "../model/profile-data";
import { ProfileHeader } from "./profile-header";
import { ProfileQrCode } from "./profile-qr-code";

export interface ProfileViewProps {
  actorId?: string;
  loadProfile: (actorId: string) => Promise<{
    success: boolean;
    data?: ActorProfile;
    message?: string;
  }>;
  ownProfile?: boolean;
}

export function ProfileView({
  actorId,
  loadProfile,
  ownProfile = false,
}: ProfileViewProps) {
  const [result, setResult] = useState<ActorProfile>();
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState(true);
  const [profileUrl, setProfileUrl] = useState(
    actorId ? `/profile/${actorId}` : "/profile",
  );

  useEffect(() => {
    const publicIdentifier = result?.handle || actorId;
    setProfileUrl(
      new URL(
        publicIdentifier ? `/profile/${publicIdentifier}` : "/profile",
        window.location.origin,
      ).toString(),
    );
  }, [actorId, result?.handle]);

  useEffect(() => {
    let active = true;
    setError(undefined);
    setResult(undefined);
    if (!actorId) {
      setError(
        "Não encontramos os dados do seu usuário. Você ainda pode editar ou configurar o perfil.",
      );
      setLoading(false);
      return () => {
        active = false;
      };
    }
    setLoading(true);
    loadProfile(actorId)
      .then((response) => {
        if (!active) return;
        if (response.success && response.data) setResult(response.data);
        else
          setError(
            response.message || "Não foi possível carregar este perfil.",
          );
      })
      .catch((cause: unknown) => {
        if (!active) return;
        setError(
          cause instanceof Error
            ? cause.message
            : "Não foi possível carregar este perfil.",
        );
      })
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [actorId, loadProfile]);

  useEffect(() => {
    const updateProfile = (event: Event) => {
      const profile = (event as CustomEvent<ActorProfile["profile"]>).detail;
      setResult((current) => (current ? { ...current, profile } : current));
    };
    window.addEventListener("univents:profile-updated", updateProfile);
    return () =>
      window.removeEventListener("univents:profile-updated", updateProfile);
  }, []);

  if (loading) return <ProfileSkeleton />;
  if (error && !ownProfile) return <MissingPublicProfile />;

  const profile = asUniventsProfile(result?.profile ?? {});
  const name = profileDisplayName(profile);
  const socials = Object.entries(profile.socials ?? {}).filter(
    (entry): entry is [string, string] => Boolean(entry[1]),
  );

  const specializations = profile.languages ?? [];

  const hasContact =
    profile.website || profile.contactEmail || socials.length > 0;
  const completeness = profileCompleteness(profile);

  return (
    <main className="min-h-dvh bg-background pb-28">
      <ProfileHeader
        profile={profile}
        name={name}
        ownProfile={ownProfile}
        profileUrl={profileUrl}
      />

      <div className="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:mt-5 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
        {/* ---- Main Column ---- */}
        <div className="space-y-5">
          {ownProfile && completeness < 100 && (
            <ProfileCard title="Integridade do Perfil">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-muted-foreground">
                  {completeness >= 100
                    ? "Perfil completo!"
                    : "Complete seu perfil"}
                </span>
                <span className="rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">
                  {completeness}% completo
                </span>
              </div>
              <div className="mt-3 h-2 w-full overflow-hidden rounded-full bg-muted">
                <div
                  className="h-full rounded-full bg-primary transition-all"
                  style={{
                    width: `${completeness}%`,
                  }}
                />
              </div>
              <p className="mt-3 text-sm italic text-muted-foreground">
                {completenessHint(profile)}
              </p>
            </ProfileCard>
          )}

          {error && ownProfile && (
            <p role="status" className="text-sm text-muted-foreground">
              Seu perfil ainda não foi criado. Use “Editar perfil” para começar.
            </p>
          )}

          <ProfileCard title="Sobre mim">
            {profile.aboutMe ? (
              <p className="whitespace-pre-wrap text-[15px] leading-[1.7] text-muted-foreground">
                {profile.aboutMe}
              </p>
            ) : (
              <EmptyState message="Nada aqui ainda. Conte um pouco sobre você!" />
            )}

            <div className="mt-6 flex flex-wrap gap-6 border-t border-border pt-5">
              {profile.createdAt ? (
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-md bg-primary/10 text-primary">
                    <Calendar className="size-4" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-card-foreground">
                      Membro desde
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {formatMemberSince(profile.createdAt)}
                    </p>
                  </div>
                </div>
              ) : (
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-md bg-muted text-muted-foreground">
                    <Calendar className="size-4" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-card-foreground">
                      Membro desde
                    </p>
                    <p className="text-sm text-muted-foreground">—</p>
                  </div>
                </div>
              )}
            </div>
          </ProfileCard>

          <ProfileCard title="Idiomas">
            {specializations.length > 0 ? (
              <div className="flex flex-wrap gap-2">
                {specializations.map((item: string) => (
                  <span
                    key={item}
                    className="rounded-md border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted"
                  >
                    {item}
                  </span>
                ))}
              </div>
            ) : (
              <EmptyState message="Adicione suas especializações para destacar seu perfil." />
            )}
          </ProfileCard>
        </div>

        {/* ---- Sidebar ---- */}
        <div className="space-y-5">
          {ownProfile && (
            <div className="hidden md:block">
              <ProfileCard title="Compartilhar perfil">
                <div className="flex flex-col items-center gap-3 text-center">
                  <ProfileQrCode value={profileUrl} />
                  <p className="text-sm text-muted-foreground">
                    Escaneie para abrir e salvar este perfil.
                  </p>
                  <button
                    type="button"
                    onClick={() => navigator.clipboard.writeText(profileUrl)}
                    className={buttonVariants({
                      variant: "outline",
                      className: "h-9 w-full rounded-md",
                    })}
                  >
                    <Copy className="mr-2 size-4" />
                    Copiar link do perfil
                  </button>
                </div>
              </ProfileCard>
            </div>
          )}

          <ProfileCard title="Contato">
            {hasContact ? (
              <div className="grid grid-cols-2 gap-1">
                {profile.website && (
                  <SocialButton
                    href={profile.website}
                    label="Website"
                    icon={<Globe className="size-4" />}
                  />
                )}
                {profile.contactEmail && (
                  <div className="hidden md:block">
                    <SocialButton
                      href={`mailto:${profile.contactEmail}`}
                      label="E-mail"
                      icon={<Mail className="size-4" />}
                    />
                  </div>
                )}
                {socials.map(([network, value]) => (
                  <SocialButton
                    key={network}
                    href={socialHref(network, value)}
                    label={capitalize(network)}
                    icon={<SocialIcon network={network} />}
                  />
                ))}
              </div>
            ) : (
              <EmptyState message="Adicione formas de contato para que outros possam se conectar." />
            )}
          </ProfileCard>
        </div>
      </div>
    </main>
  );
}

/* ========== Subcomponents ========== */

function ProfileCard({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section
      className={cn(
        "rounded-md border border-border bg-card p-5",
        "shadow-md shadow-foreground/5",
      )}
    >
      <h2 className="mb-4 text-base font-semibold text-card-foreground">
        {title}
      </h2>
      {children}
    </section>
  );
}

function SocialButton({
  href,
  label,
  icon,
}: {
  href: string;
  label: string;
  icon?: React.ReactNode;
}) {
  return (
    <a
      href={href}
      target={href.startsWith("mailto:") ? undefined : "_blank"}
      rel="noreferrer"
      className="flex items-center gap-2.5 rounded-md p-2 text-sm text-card-foreground transition-colors hover:bg-muted"
    >
      <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-foreground text-background">
        {icon}
      </span>
      <span className="truncate">{label}</span>
    </a>
  );
}

const SOCIAL_ICONS: Record<string, string> = {
  x: xIcon,
  github: githubIcon,
  twitch: twitchIcon,
  bluesky: blueskyIcon,
  discord: discordIcon,
  youtube: youtubeIcon,
  linkedin: linkedinIcon,
  instagram: instagramIcon,
};

function SocialIcon({ network }: { network: string }) {
  const src = SOCIAL_ICONS[network];
  const preserveColors = ["youtube", "bluesky", "linkedin"].includes(network);
  return src ? (
    <img
      src={src}
      alt={`${capitalize(network)} — ícone`}
      className={`size-4 object-contain${preserveColors ? "" : " dark:invert"}`}
    />
  ) : (
    <Globe className="size-4" />
  );
}

function EmptyState({ message }: { message: string }) {
  return <p className="text-sm italic text-muted-foreground/70">{message}</p>;
}

function ProfileSkeleton() {
  return (
    <main className="min-h-dvh bg-background pb-28">
      <div className="w-full">
        <Skeleton className="h-44 w-full md:h-56" />
        <div className="mx-auto max-w-7xl px-4">
          <Skeleton className="-mt-14 size-28 rounded-full border-4 border-border md:-mt-16 md:size-32" />
          <div className="mt-4 space-y-2">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-4 w-64" />
            <Skeleton className="h-4 w-40" />
          </div>
        </div>
      </div>
      <div className="mx-auto mt-5 hidden max-w-7xl gap-5 px-4 md:grid md:grid-cols-[1fr_280px]">
        <div className="space-y-5">
          <Skeleton className="h-36 rounded-md" />
          <Skeleton className="h-32 rounded-md" />
        </div>
        <div className="space-y-5">
          <Skeleton className="h-64 rounded-md" />
          <Skeleton className="h-24 rounded-md" />
        </div>
      </div>
    </main>
  );
}

function MissingPublicProfile() {
  return (
    <main className="relative min-h-dvh overflow-hidden bg-background">
      <div
        aria-hidden="true"
        className="pointer-events-none blur-md opacity-50"
      >
        <ProfileSkeleton />
      </div>
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/35 px-4 backdrop-blur-sm">
        <div
          role="alert"
          className="w-full max-w-md rounded-lg border border-border bg-card p-6 text-center shadow-xl"
        >
          <CircleAlert className="mx-auto size-8 text-muted-foreground" />
          <h1 className="mt-4 text-xl font-semibold">Perfil não encontrado</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Este perfil não existe ou não está mais disponível.
          </p>
          <Link
            to="/profile"
            className={buttonVariants({ className: "mt-5 w-full" })}
          >
            Ir para o meu perfil
          </Link>
        </div>
      </div>
    </main>
  );
}

/* ========== Helpers ========== */

function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function formatMemberSince(iso: string) {
  const d = new Date(iso);
  return d.toLocaleDateString("en-US", { month: "long", year: "numeric" });
}

function completenessHint(profile: ReturnType<typeof asUniventsProfile>) {
  if (!(profile.preferredName || profile.legalName))
    return "Adicione seu nome para completar o perfil.";
  if (!profile.pfpUrl) return "Adicione uma foto para completar o perfil.";
  if (!profile.bannerUrl)
    return "Adicione uma imagem de capa para completar o perfil.";
  if (!profile.aboutMe)
    return 'Preencha a seção "Sobre mim" para completar o perfil.';
  if (!(profile.role || profile.organization))
    return "Adicione sua função ou organização para completar o perfil.";
  if (!profile.languages?.length)
    return "Adicione ao menos um idioma para completar o perfil.";
  if (
    !(
      profile.website ||
      profile.contactEmail ||
      Object.values(profile.socials ?? {}).some(Boolean)
    )
  )
    return "Adicione uma forma de contato para completar o perfil.";
  return "Perfil completo!";
}
