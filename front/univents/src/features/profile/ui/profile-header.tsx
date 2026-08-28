import { Link } from "@tanstack/react-router";
import { Calendar, Mail, Pencil, Settings } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/shadcn/avatar";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import type { UniventsProfile } from "../model/profile-data";
import { ProfileShareDialog } from "./profile-share-dialog";

interface ProfileHeaderProps {
  profile: UniventsProfile;
  name: string;
  handle?: string;
  ownProfile: boolean;
  profileUrl: string;
  activeTab: "about" | "badges" | "purchases";
  onTabChange: (tab: "about" | "badges" | "purchases") => void;
}

export function ProfileHeader({
  profile,
  name,
  handle,
  ownProfile,
  profileUrl,
  activeTab,
  onTabChange,
}: ProfileHeaderProps) {
  const bannerStyle = profile.bannerUrl
    ? {
        backgroundImage: `url(${profile.bannerUrl})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
      }
    : undefined;

  return (
    <section className="relative w-full border-b border-border bg-card shadow-md">
      <div
        className={cn(
          "h-40 w-full bg-background bg-linear-to-br",
          "from-primary/40 via-primary/15 to-muted sm:h-48 md:h-56",
        )}
        style={bannerStyle}
      />

      <div className="md:hidden">
        {ownProfile && (
          <ProfileShareDialog
            profileUrl={profileUrl}
            className="absolute right-4 top-4"
          />
        )}
        <div className="relative mx-auto px-4">
          <div className="flex -mt-12 flex-col items-center">
            <ProfileAvatar
              name={name}
              imageUrl={profile.pfpUrl}
              size="mobile"
            />
          </div>
          <div className="mt-3 pb-4 text-center">
            <div className="flex items-baseline justify-center gap-2">
              <h1 className="text-2xl font-semibold text-card-foreground">
                {name}
              </h1>
              {profile.pronouns && (
                <span className="text-sm text-muted-foreground">
                  {profile.pronouns}
                </span>
              )}
            </div>
            {handle && (
              <p className="text-sm text-muted-foreground">
                @{handle.replace(/^@/, "")}
              </p>
            )}
            {profile.legalName &&
              profile.legalName !== name &&
              profile.visibility?.hideLegalName === false && (
                <p className="text-sm text-muted-foreground">
                  {profile.legalName}
                </p>
              )}
            {(profile.role || profile.organization) && (
              <p className="mt-1 text-sm text-muted-foreground">
                {[profile.role, profile.organization]
                  .filter(Boolean)
                  .join(" · ")}
              </p>
            )}
            <MemberSince createdAt={profile.createdAt} centered />
            {profile.contactEmail && (
              <a
                href={`mailto:${profile.contactEmail}`}
                className="mt-1 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
              >
                <Mail className="size-3.5" />
                {profile.contactEmail}
              </a>
            )}
          </div>
          {ownProfile && (
            <div className="mb-4 flex gap-2">
              <Link
                to="/profile/edit"
                search={{ returnTo: undefined }}
                className={buttonVariants({
                  className: "h-11 flex-1 rounded-lg shadow-sm",
                })}
              >
                <Pencil className="mr-2 size-4" />
                Editar perfil
              </Link>
              <Link
                to="/profile/config"
                className={buttonVariants({
                  variant: "outline",
                  size: "icon",
                  className: "size-11 shrink-0 rounded-lg shadow-sm",
                })}
                aria-label="Configurações do perfil"
              >
                <Settings className="size-4" />
              </Link>
            </div>
          )}
        </div>
      </div>

      <div className="mx-auto hidden max-w-7xl px-4 md:block">
        <div className="relative pb-6">
          <div className="absolute -top-16">
            <ProfileAvatar
              name={name}
              imageUrl={profile.pfpUrl}
              size="desktop"
            />
          </div>
          {ownProfile && <DesktopActions profileUrl={profileUrl} />}
          <div className="pt-18">
            <div className="flex items-baseline gap-2">
              <h1 className="text-3xl font-semibold tracking-tight text-card-foreground">
                {name}
              </h1>
              {profile.pronouns && (
                <span className="text-sm text-muted-foreground">
                  {profile.pronouns}
                </span>
              )}
            </div>
            {handle && (
              <p className="text-sm text-muted-foreground">
                @{handle.replace(/^@/, "")}
              </p>
            )}
            {profile.legalName &&
              profile.legalName !== name &&
              profile.visibility?.hideLegalName === false && (
                <p className="text-sm text-muted-foreground">
                  {profile.legalName}
                </p>
              )}
            {(profile.role || profile.organization) && (
              <p className="text-[15px] text-muted-foreground">
                {[profile.role, profile.organization]
                  .filter(Boolean)
                  .join(" · ")}
              </p>
            )}
            <MemberSince createdAt={profile.createdAt} />
          </div>
        </div>
      </div>
      <div className="relative mx-auto max-w-7xl">
        <nav className="flex overflow-x-auto px-4 pr-12 sm:pr-4">
          {(
            [
              ["about", "Sobre"],
              ["badges", "Crachás"],
              ...(ownProfile ? ([["purchases", "Compras"]] as const) : []),
            ] as Array<["about" | "badges" | "purchases", string]>
          ).map(([tab, label]) => (
            <button
              key={tab}
              type="button"
              onClick={() => onTabChange(tab)}
              className={cn(
                "shrink-0 whitespace-nowrap border-b-2 px-4 py-3 text-sm font-medium transition-colors",
                activeTab === tab
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground",
              )}
            >
              {label}
            </button>
          ))}
        </nav>
        {ownProfile ? (
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-y-0 right-0 w-12 bg-linear-to-r from-transparent to-card sm:hidden"
          />
        ) : null}
      </div>
    </section>
  );
}

function MemberSince({
  createdAt,
  centered = false,
}: {
  createdAt?: string;
  centered?: boolean;
}) {
  return (
    <p
      className={cn(
        "mt-1 flex items-center gap-1.5 text-sm text-muted-foreground",
        centered && "justify-center",
      )}
    >
      <Calendar className="size-3.5" />
      Membro desde {createdAt ? formatMemberSince(createdAt) : "—"}
    </p>
  );
}

function formatMemberSince(iso: string) {
  return new Date(iso).toLocaleDateString("pt-BR", {
    month: "long",
    year: "numeric",
  });
}

function ProfileAvatar({
  name,
  imageUrl,
  size,
}: {
  name: string;
  imageUrl?: string | null;
  size: "mobile" | "desktop";
}) {
  return (
    <Avatar
      className={cn(
        "border-4 border-background bg-background shadow-xl",
        size === "mobile" ? "size-24" : "size-32",
      )}
    >
      <AvatarImage src={imageUrl ?? undefined} alt={name} />
      <AvatarFallback className="text-2xl font-semibold">
        {name.slice(0, 2).toUpperCase()}
      </AvatarFallback>
    </Avatar>
  );
}

function DesktopActions({ profileUrl }: { profileUrl: string }) {
  return (
    <div className="absolute right-0 top-4 flex gap-2">
      <ProfileShareDialog
        profileUrl={profileUrl}
        className="size-10 rounded-md shadow-sm"
      />
      <Link
        to="/profile/edit"
        search={{ returnTo: undefined }}
        className={buttonVariants({
          className: "h-10 rounded-md px-4 shadow-sm",
        })}
      >
        <Pencil className="size-4" />
        Editar perfil
      </Link>
      <Link
        to="/profile/config"
        className={buttonVariants({
          variant: "outline",
          size: "icon",
          className: "size-10 rounded-md shadow-sm",
        })}
        aria-label="Configurações"
      >
        <Settings className="size-4" />
      </Link>
    </div>
  );
}
