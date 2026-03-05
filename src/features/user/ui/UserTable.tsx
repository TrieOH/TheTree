import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import type { User } from "../model/types";
import { formatDate } from "@/shared/lib/date-utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import UserPermEditor from "./UserPermEditor";
import { User as UserIcon } from "lucide-react";
import TruncatedId from "@/shared/ui/TruncatedId";

interface PropsI {
  data: User[];
  project_id: string;
}

const UserInfo = ({ user }: { user: User }) => {
  return (
    <div className="flex items-center gap-3 py-1 text-left">
      <div className="shrink-0 w-10 h-10 rounded-full bg-muted flex items-center justify-center border border-border">
        <UserIcon className="w-5 h-5 text-muted-foreground" />
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-base font-bold text-foreground truncate">{user.email}</span>
          {user.is_verified && (
            <Badge variant="default" className="bg-emerald-500/10 text-emerald-500 border-emerald-500/20 text-[10px] px-1.5 py-0 h-4 uppercase tracking-wider font-bold">Verified</Badge>
          )}
        </div>
        <div className="text-xs text-muted-foreground mt-0.5">
          <TruncatedId id={user.id} />
        </div>
      </div>
    </div>
  );
};

export default function UserTable({ data, project_id }: PropsI) {
  return (
    <CustomDataTable
      forceMobileView={true}
      mobileColumnCount={5}
      data={data}
      columns={[
        {
          key: "email",
          header: "User",
          primary: true,
          sortable: true,
          render: (_, row) => <UserInfo user={row} />,
        },
        {
          key: "created_at",
          header: "Joined At",
          sortable: true,
          render: (value) => formatDate(value as string),
          searchableTextExtractor: (value) => formatDate(value as string),
        },
        {
          key: "last_login_at",
          header: "Last Login",
          sortable: true,
          render: (value) => (value ? formatDate(value as string) : "Never"),
          searchableTextExtractor: (value) => (value ? formatDate(value as string) : "Never"),
        },
        {
          key: "verified_at",
          header: "Verified At",
          sortable: true,
          render: (value) => (value ? formatDate(value as string) : "Unverified"),
          searchableTextExtractor: (value) => (value ? formatDate(value as string) : "Unverified"),
        },
      ]}
      renderExpandedRow={(row) => <UserPermEditor project_id={project_id} user={row} />}
    />
  )
}
