import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { format } from "date-fns";
import { ptBR } from "date-fns/locale";
import {
  CalendarDays,
  CalendarPlus,
  CheckCircle2,
  ChevronRight,
  CircleAlert,
  CreditCard,
  Eye,
  LayoutGrid,
  Users,
} from "lucide-react";
import { motion } from "motion/react";
import { useState } from "react";
import { toast } from "sonner";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import {
  useDiscontinueEventMutation,
  usePublishEventMutation,
} from "@/features/events/api/mutations";
import {
  useConnectEventSellerMutation,
  useDisconnectEventSellerMutation,
} from "@/features/payments/api/mutations";
import type { PaymentProviderI } from "@/features/payments/model";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { QuickAction } from "@/widgets/ui/quick-action";

const statusConfig = {
  draft: {
    label: "Rascunho",
    className: "bg-amber-500/10 text-amber-700 border-amber-500/20",
    dot: "bg-amber-500",
  },
  active: {
    label: "Ativo",
    className: "bg-emerald-500/10 text-emerald-700 border-emerald-500/20",
    dot: "bg-emerald-500",
  },
  archived: {
    label: "Arquivado",
    className: "bg-slate-500/10 text-slate-700 border-slate-500/20",
    dot: "bg-slate-500",
  },
  discontinued: {
    label: "Descontinuado",
    className: "bg-rose-500/10 text-rose-700 border-rose-500/20",
    dot: "bg-rose-500",
  },
} as const;

const paymentProviders: Array<{
  id: PaymentProviderI;
  name: string;
  description: string;
}> = [
  {
    id: "mercadopago",
    name: "Mercado Pago",
    description: "Receba por Pix e cartão de crédito.",
  },
];

export const Route = createLazyFileRoute("/admin/events/$eventId/")({
  component: EventOverviewRoute,
});

function EventOverviewRoute() {
  const { eventId } = Route.useParams();
  const { data: ownedEvents = [] } = useQuery(allOwnEventsQueryOptions());
  const { data: joinedEvents = [] } = useQuery(allJoinedEventsQueryOptions());
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false);
  const [discontinueConfirmOpen, setDiscontinueConfirmOpen] = useState(false);
  const [disconnectSellerConfirmOpen, setDisconnectSellerConfirmOpen] =
    useState(false);
  const event =
    [...ownedEvents, ...joinedEvents].find((item) => item.id === eventId) ??
    null;
  const isPublished = event?.status === "active";
  const status = event ? statusConfig[event.status] : statusConfig.draft;

  const publishEventMutation = usePublishEventMutation();
  const discontinueEventMutation = useDiscontinueEventMutation();
  const connectSellerMutation = useConnectEventSellerMutation();
  const disconnectSellerMutation = useDisconnectEventSellerMutation();

  const copyLink = () => {
    if (!event) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${event.slug}`,
    );
    toast.success("Link copiado");
  };

  const handlePublishEvent = () => {
    if (!event || isPublished) return;
    publishEventMutation.mutate(eventId);
  };

  const handleDiscontinueEvent = () => {
    if (!event || event.status !== "active") return;
    discontinueEventMutation.mutate(eventId);
  };

  const metrics = [
    {
      label: "Criado em",
      value: event
        ? format(new Date(event.created_at), "dd MMM yyyy", { locale: ptBR })
        : "—",
      hint: "Data de criação do evento",
    },
    {
      label: "Atualizado em",
      value: event?.updated_at
        ? format(new Date(event.updated_at), "dd MMM yyyy", { locale: ptBR })
        : "—",
      hint: "Última alteração registrada",
    },
    {
      label: "Contato",
      value: event?.contact_email ?? "—",
      hint: "E-mail principal do evento",
    },
  ];
  const heroDescription =
    event?.description ?? "Sem descrição cadastrada para este evento.";

  const checklist = [
    {
      label: "Banner e logo cadastrados",
      done: Boolean(event?.banner_url || event?.logo_url),
    },
    {
      label: "Descrição preenchida",
      done: Boolean(event?.description),
    },
    {
      label: "Slug público disponível",
      done: Boolean(event?.slug),
    },
  ];

  const sections = [
    {
      label: "Edições",
      to: "/admin/events/$eventId/editions",
      icon: CalendarPlus,
    },
    {
      label: "Membros",
      to: "/admin/events/$eventId/members",
      icon: Users,
    },
  ];

  const actions = [
    ...(event?.status === "draft"
      ? [
          {
            label: "Publicar evento",
            onClick: () => setPublishConfirmOpen(true),
            disabled: publishEventMutation.isPending,
            variant: "default" as const,
          },
        ]
      : []),
    {
      label: "Copiar link público",
      onClick: copyLink,
      disabled: !event,
      variant: "default" as const,
    },
    ...(isPublished
      ? [
          {
            label: "Descontinuar evento",
            onClick: () => setDiscontinueConfirmOpen(true),
            disabled: discontinueEventMutation.isPending,
            variant: "destructive" as const,
          },
          {
            label: "Abrir painel público",
            to: "/events/$slug" as const,
            params: { slug: event?.slug ?? "" },
            variant: "default" as const,
          },
        ]
      : []),
  ];

  return (
    <>
      <div className="relative space-y-6 p-6 pb-28!">
        <motion.section
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35 }}
          className="overflow-hidden rounded-md border border-border/60 bg-card shadow-[0_1px_0_0_rgba(255,255,255,0.03),0_20px_40px_-24px_rgba(15,23,42,0.24)]"
        >
          <div className="relative flex flex-col gap-6 p-6">
            <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-primary/20 to-transparent" />
            <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-primary/5 via-transparent to-transparent" />

            <div className="space-y-5">
              <div className="flex items-center gap-2">
                <Badge variant="secondary" className="rounded-full px-3">
                  <LayoutGrid className="size-3.5" />
                  Overview
                </Badge>
              </div>

              <div className="space-y-3">
                <h1 className="max-w-3xl text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
                  {event?.full_name ?? "Evento"}
                </h1>
                <p className="max-w-2xl text-sm leading-6 text-muted-foreground md:text-base">
                  {heroDescription}
                </p>
              </div>

              <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                  <CalendarDays className="size-3.5" />
                  {event?.slug ?? "slug-do-evento"}
                </span>
                <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                  <Eye className="size-3.5" />
                  Link público
                </span>
                <span
                  className={cn(
                    "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1",
                    status.className,
                  )}
                >
                  <span className={cn("size-1.5 rounded-full", status.dot)} />
                  {status.label}
                </span>
              </div>
            </div>
          </div>
        </motion.section>

        <section className="grid gap-4 md:grid-cols-3">
          {metrics.map((metric) => (
            <Card
              key={metric.label}
              className="border-border/60 bg-card/95 shadow-sm transition-shadow hover:shadow-md"
            >
              <CardHeader className="pb-2">
                <CardDescription className="flex items-center gap-2 text-xs uppercase tracking-[0.22em]">
                  <span className="size-1.5 rounded-full bg-primary/60" />
                  {metric.label}
                </CardDescription>
                <CardTitle className="text-2xl font-semibold tracking-tight">
                  {metric.value}
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <p className="text-xs text-muted-foreground">{metric.hint}</p>
              </CardContent>
            </Card>
          ))}
        </section>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Pagamentos</CardTitle>
            <CardDescription>
              Escolha a conta que receberá as vendas deste evento.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-3 py-5 md:grid-cols-2 xl:grid-cols-3">
            {paymentProviders.map((provider) => {
              const connected = Boolean(event?.payssage_seller_id);

              return (
                <div
                  key={provider.id}
                  className="flex min-h-36 flex-col justify-between gap-5 rounded-xl border border-border/60 bg-muted/15 p-4"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex items-center gap-3">
                      <div className="rounded-lg bg-background p-2.5 shadow-sm">
                        <CreditCard className="size-5 text-muted-foreground" />
                      </div>
                      <div>
                        <p className="text-sm font-semibold">{provider.name}</p>
                        <p className="mt-1 text-xs leading-5 text-muted-foreground">
                          {provider.description}
                        </p>
                      </div>
                    </div>
                    <Badge variant={connected ? "default" : "secondary"}>
                      {connected ? "Conectado" : "Disponível"}
                    </Badge>
                  </div>

                  <Button
                    className="w-full"
                    variant={connected ? "outline" : "default"}
                    disabled={
                      !event ||
                      connectSellerMutation.isPending ||
                      disconnectSellerMutation.isPending
                    }
                    onClick={() =>
                      connected
                        ? setDisconnectSellerConfirmOpen(true)
                        : connectSellerMutation.mutate({
                            eventId,
                            provider: provider.id,
                          })
                    }
                  >
                    {connected ? "Desconectar conta" : "Conectar conta"}
                  </Button>
                </div>
              );
            })}
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Checklist do evento</CardTitle>
            <CardDescription>
              Itens derivados dos dados já cadastrados no evento.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 py-5">
            {checklist.map((item) => (
              <div
                key={item.label}
                className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/15 px-4 py-3.5"
              >
                <div className="flex items-center gap-3">
                  <div
                    className={cn(
                      "size-2 rounded-full",
                      item.done ? "bg-emerald-500" : "bg-amber-500",
                    )}
                  />
                  <span className="text-sm text-foreground">{item.label}</span>
                </div>
                {item.done ? (
                  <CheckCircle2 className="size-4 text-emerald-500/70" />
                ) : (
                  <CircleAlert className="size-4 text-amber-500/70" />
                )}
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Seções do evento</CardTitle>
            <CardDescription>
              Gerencie os recursos vinculados a este evento.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-3 py-5 sm:grid-cols-2 xl:grid-cols-4">
            {sections.map((section) => {
              const Icon = section.icon;
              return (
                <QuickAction
                  key={section.to}
                  to={section.to}
                  params={{ eventId }}
                >
                  <div className="flex items-center gap-3">
                    <Icon className="size-4 text-muted-foreground" />
                    <span className="text-sm font-medium text-foreground">
                      {section.label}
                    </span>
                  </div>
                  <ChevronRight className="size-4 text-muted-foreground" />
                </QuickAction>
              );
            })}
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle>Ações rápidas</CardTitle>
            <CardDescription>
              Atalhos para as operações mais comuns do evento.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-3 py-5 sm:grid-cols-2 xl:grid-cols-3">
            {actions.map((action) => {
              if ("to" in action && action.to) {
                return (
                  <QuickAction
                    key={action.label}
                    to={action.to}
                    params={action.params}
                    disabled={action.disabled}
                    variant={action.variant}
                  >
                    <span className="text-sm font-medium text-foreground">
                      {action.label}
                    </span>
                    <ChevronRight className="size-4 text-muted-foreground" />
                  </QuickAction>
                );
              }

              return (
                <QuickAction
                  key={action.label}
                  onClick={action.onClick}
                  disabled={action.disabled}
                  variant={action.variant}
                >
                  <span className="text-sm font-medium text-foreground">
                    {action.label}
                  </span>
                  <ChevronRight className="size-4 text-muted-foreground" />
                </QuickAction>
              );
            })}
          </CardContent>
        </Card>
      </div>

      <AlertModal
        open={disconnectSellerConfirmOpen}
        onOpenChange={setDisconnectSellerConfirmOpen}
        title="Desconectar Mercado Pago?"
        description="Este evento deixará de receber novos pagamentos até uma conta ser conectada novamente."
        confirmLabel="Desconectar"
        variant="destructive"
        loading={disconnectSellerMutation.isPending}
        onConfirm={() => {
          disconnectSellerMutation.mutate(eventId, {
            onSuccess: () => setDisconnectSellerConfirmOpen(false),
          });
        }}
      />

      <AlertModal
        open={publishConfirmOpen}
        onOpenChange={setPublishConfirmOpen}
        title="Publicar evento?"
        description="Depois de publicar, o painel público ficará disponível para o evento."
        confirmLabel="Publicar evento"
        variant="default"
        loading={publishEventMutation.isPending}
        onConfirm={async () => {
          handlePublishEvent();
          setPublishConfirmOpen(false);
        }}
      />

      <AlertModal
        open={discontinueConfirmOpen}
        onOpenChange={setDiscontinueConfirmOpen}
        title="Descontinuar evento?"
        description="O evento deixará de ser ativo e a data de atualização será atualizada."
        confirmLabel="Descontinuar evento"
        variant="destructive"
        loading={discontinueEventMutation.isPending}
        onConfirm={async () => {
          handleDiscontinueEvent();
          setDiscontinueConfirmOpen(false);
        }}
      />
    </>
  );
}
