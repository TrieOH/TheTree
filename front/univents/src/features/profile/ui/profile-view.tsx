import { useQueries, useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { queryError } from "@trieoh/front-core";
import type { ActorProfile } from "@trieoh/identityx-sdk-ts";
import { Globe, Mail } from "lucide-react";
import { userBadgesQueryOptions } from "@/features/badges/api";
import {
  type BadgeProfileBadge,
  badgeDesignSchema,
} from "@/features/badges/model";
import { allProfileBadges } from "@/features/badges/model/profile-badges";
import { BadgePreview } from "@/features/badges/ui/badge-preview";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import { PurchasesContent } from "@/features/purchases/ui/purchases-content";
import { cn } from "@/shared/lib/utils";
import blueskyIcon from "@/shared/ui/social-icons/assets/bluesky.svg";
import discordIcon from "@/shared/ui/social-icons/assets/discord.svg";
import githubIcon from "@/shared/ui/social-icons/assets/github.svg";
import instagramIcon from "@/shared/ui/social-icons/assets/instagram.svg";
import linkedinIcon from "@/shared/ui/social-icons/assets/linkedin.svg";
import twitchIcon from "@/shared/ui/social-icons/assets/twitch.svg";
import xIcon from "@/shared/ui/social-icons/assets/x.svg";
import youtubeIcon from "@/shared/ui/social-icons/assets/youtube.svg";
import { profileKeys } from "../api/query-keys";
import {
  asUniventsProfile,
  profileCompleteness,
  profileDisplayName,
  socialHref,
} from "../model/profile-data";
import { ProfileHeader } from "./profile-header";
import {
  completenessHint,
  MissingPublicProfile,
  ProfileSkeleton,
} from "./profile-view-state";

export interface ProfileViewProps {
  actorId?: string;
  loadProfile: (actorId: string) => Promise<{
    success: boolean;
    data?: ActorProfile;
    message?: string;
    code?: number;
  }>;
  ownProfile?: boolean;
  viewerActorId?: string;
  activeTab: "about" | "badges" | "certificates" | "purchases";
  onTabChange: (tab: "about" | "badges" | "certificates" | "purchases") => void;
}

export function ProfileView({
  actorId,
  loadProfile,
  ownProfile = false,
  viewerActorId,
  activeTab,
  onTabChange,
}: ProfileViewProps) {
  const profileQuery = useQuery({
    queryKey: profileKeys.detail(actorId ?? ""),
    enabled: Boolean(actorId),
    queryFn: async () => {
      const response = await loadProfile(actorId ?? "");
      if (!response.success || !response.data) {
        throw queryError(
          response.message || "Não foi possível carregar este perfil.",
          response.code,
        );
      }
      return response.data;
    },
  });
  const result = profileQuery.data;
  const { data: badges } = useQuery(
    userBadgesQueryOptions(result?.actor_id ?? ""),
  );
  const { data: events = [] } = useQuery(allPublicEventsQueryOptions());
  const editionQueries = useQueries({
    queries: events.map((event) => allPublicEditionsQueryOptions(event.id)),
  });
  const editionLocations = new Map(
    editionQueries
      .flatMap((query) => query.data ?? [])
      .map((edition) => [edition.id, edition.location_name ?? ""]),
  );
  const error = actorId
    ? profileQuery.error?.message
    : "Não encontramos os dados do seu usuário. Você ainda pode editar ou configurar o perfil.";

  if (profileQuery.isLoading) return <ProfileSkeleton />;
  if (error && !ownProfile) return <MissingPublicProfile />;

  const profile = asUniventsProfile({
    ...(result?.profile ?? {}),
    ...(result?.pfp_url !== undefined ? { pfpUrl: result.pfp_url } : {}),
  });
  const isOwnProfile = ownProfile || result?.actor_id === viewerActorId;
  const publicIdentifier = result?.handle || actorId;
  const profileUrl = new URL(
    publicIdentifier ? `/profile/${publicIdentifier}` : "/",
    window.location.origin,
  ).toString();
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
        handle={result?.handle ?? undefined}
        ownProfile={isOwnProfile}
        profileUrl={profileUrl}
        activeTab={activeTab}
        onTabChange={onTabChange}
      />

      {activeTab === "badges" ? (
        <div className="mx-auto mt-5 max-w-7xl px-4">
          <ProfileBadges
            badges={badges ? allProfileBadges(badges) : []}
            profileIdentifier={publicIdentifier ?? ""}
            participantName={profile.preferredName || profile.legalName || ""}
            editionLocations={editionLocations}
          />
        </div>
      ) : activeTab === "certificates" && isOwnProfile ? (
        <div className="mx-auto mt-5 max-w-7xl px-4">
          <UserCertificationsSection
            participantName={profile.preferredName || profile.legalName || ""}
          />
        </div>
      ) : activeTab === "purchases" ? (
        <div className="mx-auto mt-4 max-w-7xl px-4">
          <PurchasesContent />
        </div>
      ) : (
        <div className="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:mt-5 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
          {/* ---- Main Column ---- */}
          <div className="space-y-5">
            {isOwnProfile && completeness < 100 && (
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

            <ProfileCard title="Sobre mim">
              {profile.aboutMe ? (
                <p className="whitespace-pre-wrap text-[15px] leading-[1.7] text-muted-foreground">
                  {profile.aboutMe}
                </p>
              ) : (
                <EmptyState message="Nada aqui ainda. Conte um pouco sobre você!" />
              )}
            </ProfileCard>
          </div>

          {/* ---- Sidebar ---- */}
          <div className="space-y-5">
            <ProfileCard title="Idiomas">
              {specializations.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {specializations.map((item: string, index) => (
                    <span
                      // biome-ignore lint/suspicious/noArrayIndexKey: duplicate languages are valid and these tags hold no state
                      key={`${item}-${index}`}
                      className="rounded-md border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted"
                    >
                      {item}
                    </span>
                  ))}
                </div>
              ) : (
                <EmptyState message="Adicione idiomas para destacar seu perfil." />
              )}
            </ProfileCard>

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
      )}
    </main>
  );
}

function ProfileBadges({
  badges,
  profileIdentifier,
  participantName,
  editionLocations,
}: {
  badges: BadgeProfileBadge[];
  profileIdentifier: string;
  participantName: string;
  editionLocations: Map<string, string>;
}) {
  if (badges.length === 0) return null;
  return (
    <div className="flex flex-wrap items-start justify-center gap-2 sm:justify-start!">
      {badges.map((badge) =>
        (() => {
          const design = badgeDesignSchema.safeParse(badge.design_data);
          const canvas = design.success
            ? design.data.canvas
            : { width: 321, height: 204 };
          const width = (160 * canvas.width) / canvas.height;
          return (
            <Link
              key={badge.emission_id}
              to="/profile/$actorId/badges/$badgeId"
              params={{
                actorId: profileIdentifier,
                badgeId: badge.emission_id,
              }}
              className="flex w-fit max-w-full items-start"
              aria-label={`Abrir badge ${badge.template_name ?? badge.edition_name}`}
            >
              <BadgePreview
                badge={badge}
                framed={false}
                participantName={participantName}
                location={editionLocations.get(badge.edition_id) ?? ""}
                className="relative h-auto max-w-full"
                style={{ width, height: "auto", maxWidth: "100%" }}
              />
            </Link>
          );
        })(),
      )}
    </div>
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
      width={16}
      height={16}
      className={`size-4 object-contain${preserveColors ? "" : " dark:invert"}`}
    />
  ) : (
    <Globe className="size-4" />
  );
}

function EmptyState({ message }: { message: string }) {
  return <p className="text-sm italic text-muted-foreground/70">{message}</p>;
}

/* ========== Helpers ========== */

function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}
