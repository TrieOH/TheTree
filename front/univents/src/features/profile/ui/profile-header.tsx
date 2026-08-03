import { Link } from "@tanstack/react-router";
import { Pencil, Settings } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/shadcn/avatar";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import type { UniventsProfile } from "../model/profile-data";
import { ProfileShareDialog } from "./profile-share-dialog";

interface ProfileHeaderProps {
  profile: UniventsProfile;
  name: string;
  ownProfile: boolean;
  profileUrl: string;
}

export function ProfileHeader({
  profile,
  name,
  ownProfile,
  profileUrl,
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
          "h-40 w-full bg-linear-to-br",
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
            <h1 className="text-2xl font-semibold text-card-foreground">
              {name}
            </h1>
            {(profile.role || profile.organization) && (
              <p className="mt-1 text-sm text-muted-foreground">
                {[
                  profile.role,
                  !profile.visibility?.hideOrganization
                    ? profile.organization
                    : undefined,
                ]
                  .filter(Boolean)
                  .join(" · ")}
              </p>
            )}
            {profile.pronouns && (
              <p className="text-sm text-muted-foreground">
                {profile.pronouns}
              </p>
            )}
          </div>
          {ownProfile && (
            <div className="mb-4 flex gap-2">
              <Link
                to="/profile/edit"
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
          {ownProfile && <DesktopActions />}
          <div className="pt-18">
            <h1 className="text-3xl font-semibold tracking-tight text-card-foreground">
              {name}
            </h1>
            {(profile.role || profile.organization) && (
              <p className="text-[15px] text-muted-foreground">
                {[
                  profile.role,
                  !profile.visibility?.hideOrganization
                    ? profile.organization
                    : undefined,
                ]
                  .filter(Boolean)
                  .join(" · ")}
              </p>
            )}
            {profile.pronouns && (
              <p className="text-sm text-muted-foreground">
                {profile.pronouns}
              </p>
            )}
          </div>
        </div>
      </div>
    </section>
  );
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
        "border-4 border-background shadow-xl",
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

function DesktopActions() {
  return (
    <div className="absolute right-0 top-4 flex gap-2">
      <Link
        to="/profile/edit"
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
