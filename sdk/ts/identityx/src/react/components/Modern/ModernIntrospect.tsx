import { useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { Copy, Check, User, FolderOpen, Key, Mail } from "lucide-react";

type ActorType = "human" | "service" | "machine";
type CredentialType = "token" | "api_key";

export interface IntrospectData {
  cred: {
    id?: string;
    type?: CredentialType;
  };
  sub: {
    capabilities?: number[];
    email?: string;
    id?: string;
    metadata?: number[];
    project_id?: string;
    type?: ActorType;
  };
}

export interface IdentityCardProps {
  data: IntrospectData;
}

const CAPABILITY_LABELS: Record<number, string> = {
  0: "...",
  1: "...",
  2: "...",
  3: "...",
  4: "...",
};

const ACTOR_THEME: Record<ActorType, {
  gradient: string;
  orbBg: string;
  orbText: string;
  pillBg: string;
  pillText: string;
  pillBorder: string;
  capBg: string;
  capText: string;
  capBorder: string;
}> = {
  human: {
    gradient: "linear-gradient(135deg, #7F77DD, #AFA9EC)",
    orbBg: "#EEEDFE", orbText: "#534AB7",
    pillBg: "#EEEDFE", pillText: "#3C3489", pillBorder: "#AFA9EC",
    capBg: "#EEEDFE", capText: "#534AB7", capBorder: "#AFA9EC",
  },
  service: {
    gradient: "linear-gradient(135deg, #1D9E75, #5DCAA5)",
    orbBg: "#E1F5EE", orbText: "#0F6E56",
    pillBg: "#E1F5EE", pillText: "#085041", pillBorder: "#5DCAA5",
    capBg: "#E1F5EE", capText: "#0F6E56", capBorder: "#5DCAA5",
  },
  machine: {
    gradient: "linear-gradient(135deg, #BA7517, #EF9F27)",
    orbBg: "#FAEEDA", orbText: "#854F0B",
    pillBg: "#FAEEDA", pillText: "#412402", pillBorder: "#EF9F27",
    capBg: "#FAEEDA", capText: "#854F0B", capBorder: "#EF9F27",
  },
};

function CopyIconButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    if (!value || value === "—") return;
    navigator.clipboard.writeText(value).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <button
      onClick={handleCopy}
      aria-label="Copiar valor"
      className="opacity-0 group-hover:opacity-100 transition-all duration-150 p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground"
    >
      <AnimatePresence mode="wait" initial={false}>
        {copied ? (
          <motion.span
            key="check"
            initial={{ scale: 0.7, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.7, opacity: 0 }}
            transition={{ duration: 0.15 }}
          >
            <Check className="w-3.5 h-3.5 text-green-600" />
          </motion.span>
        ) : (
          <motion.span
            key="copy"
            initial={{ scale: 0.7, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.7, opacity: 0 }}
            transition={{ duration: 0.15 }}
          >
            <Copy className="w-3.5 h-3.5" />
          </motion.span>
        )}
      </AnimatePresence>
    </button>
  );
}

function FieldRow({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ElementType;
  label: string;
  value?: string;
}) {
  const display = value || "—";
  return (
    <div className="group flex items-center justify-between px-5 py-2.5 gap-3 rounded-md mx-1.5 hover:bg-muted/50 transition-colors">
      <div className="flex items-center gap-2.5 min-w-0">
        <Icon className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
        <span className="text-[11px] text-muted-foreground w-12 shrink-0">{label}</span>
        <span
          className="text-[12px] font-mono text-foreground truncate"
          title={display}
        >
          {display}
        </span>
      </div>
      <CopyIconButton value={display} />
    </div>
  );
}

export function ModernIntrospect({ data }: IdentityCardProps) {
  const { sub, cred } = data;
  const actorType: ActorType = sub.type ?? "human";
  const theme = ACTOR_THEME[actorType];
  const caps = sub.capabilities ?? [];
  const capLabels = caps.map((c) => CAPABILITY_LABELS[c] ?? String(c));

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
      className="w-full rounded-2xl border border-border/60 overflow-hidden bg-card shadow-sm"
    >
      {/* Header */}
      <div className="px-5 pt-2.5 pb-3.5 border-b border-border/60">
        <p className="text-base font-semibold text-foreground truncate">
          {sub.id ?? "unknown_actor"}
        </p>
        <p className="text-[12px] text-muted-foreground mb-2 truncate">
          {sub.email ?? "—"}
        </p>
        <div className="flex gap-1.5 flex-wrap">
          <span
            className="text-[10px] font-medium uppercase tracking-wide px-2 py-0.5 rounded-full border"
            style={{
              background: theme.pillBg,
              color: theme.pillText,
              borderColor: theme.pillBorder,
            }}
          >
            {actorType}
          </span>
          <span className="text-[10px] font-medium uppercase tracking-wide px-2 py-0.5 rounded-full border border-border bg-muted text-muted-foreground">
            {cred.type ?? "token"}
          </span>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 border-b border-border/60 divide-x divide-border/60">
        {[
          { label: "Capabilities", value: caps.length },
          { label: "Metadata", value: (sub.metadata ?? []).length },
          { label: "Project", value: sub.project_id?.split("_")[1]?.slice(0, 4) ?? "—" },
        ].map((s) => (
          <div key={s.label} className="py-3 text-center">
            <p className="text-[10px] text-muted-foreground uppercase tracking-wide mb-0.5">{s.label}</p>
            <p className="text-sm font-medium font-mono text-foreground">{s.value}</p>
          </div>
        ))}
      </div>

      {/* Fields */}
      <div className="py-1.5 border-b border-border/60">
        <FieldRow icon={User} label="sub id" value={sub.id} />
        <FieldRow icon={FolderOpen} label="project" value={sub.project_id} />
        <FieldRow icon={Key} label="cred id" value={cred.id} />
        <FieldRow icon={Mail} label="email" value={sub.email} />
      </div>

      {/* Capabilities */}
      {capLabels.length > 0 && (
        <div className="px-5 py-3 border-b border-border/60">
          <p className="text-[10px] text-muted-foreground uppercase tracking-wide mb-2">Capabilities</p>
          <div className="flex flex-wrap gap-1.5">
            {capLabels.map((cap) => (
              <span
                key={cap}
                className="text-[11px] font-mono px-2.5 py-1 rounded-full border"
                style={{
                  background: theme.capBg,
                  color: theme.capText,
                  borderColor: theme.capBorder,
                }}
              >
                {cap}
              </span>
            ))}
          </div>
        </div>
      )}
    </motion.div>
  );
}
