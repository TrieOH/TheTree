import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import type { User } from "../model/types";
import { formatDate } from "@/shared/lib/date-utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import TruncatedId from "@/shared/ui/TruncatedId";


interface PropsI {
  data: User[]
}

export default function UserTable({ data }: PropsI) {
  return (
    <CustomDataTable
      data={data}
      columns={[
        {
          key: "email",
          header: "Email",
          sortable: true,
        },
        {
          key: "is_active",
          header: "Active",
          sortable: true,
          render: (value) =>
            value ? (
              <Badge variant="default">Yes</Badge>
            ) : (
              <Badge variant="destructive">No</Badge>
            ),
        },
        {
          key: "is_verified",
          header: "Verified",
          sortable: true,
          render: (value) =>
            value ? (
              <Badge variant="default">Yes</Badge>
            ) : (
              <Badge variant="destructive">No</Badge>
            ),
        },
        {
          key: "user_type",
          header: "Type",
          sortable: true,
        },
        {
          key: "id",
          header: "ID",
          sortable: true,
          render: (value) => <TruncatedId id={value as string} />,
        },
        {
          key: "created_at",
          header: "Created At",
          sortable: true,
          render: (value) => formatDate(value as string),
          searchableTextExtractor: (value) => formatDate(value as string),
        },
        {
          key: "updated_at",
          header: "Updated At",
          sortable: true,
          render: (value) => formatDate(value as string),
          searchableTextExtractor: (value) => formatDate(value as string),
        },
        {
          key: "last_login_at",
          header: "Last Login",
          sortable: true,
          render: (value) => formatDate(value as string),
          searchableTextExtractor: (value) => formatDate(value as string),
        },
        {
          key: "verified_at",
          header: "Verified At",
          sortable: true,
          render: (value) => formatDate(value as string),
          searchableTextExtractor: (value) => formatDate(value as string),
        },
      ]}
    />
  )
}