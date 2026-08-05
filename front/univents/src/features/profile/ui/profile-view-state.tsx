import { Link } from "@tanstack/react-router";
import { CircleAlert } from "lucide-react";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import type { UniventsProfile } from "../model/profile-data";

export function ProfileSkeleton() {
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

export function MissingPublicProfile() {
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

export function completenessHint(profile: UniventsProfile) {
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
