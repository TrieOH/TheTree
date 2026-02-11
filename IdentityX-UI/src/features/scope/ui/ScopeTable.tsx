import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Globe } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import { Scope } from "@/features/scope/model/types";

interface PropsI {
  data: Scope[]
}

export default function ScopeTable({ data }: PropsI) {
  return (
    <CustomDataTable
      data={data}
      columns={[
        {
          key: "name",
          header: "Name",
          sortable: true,
        },
        {
          key: "type",
          header: "Type",
          sortable: true,
        },
        {
          key: "external_id",
          header: "External ID",
          sortable: true,
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
      ]}
      tableActions={[
        {
          label: "Create Scope",
          icon: Globe,
          onClick: () => console.log("Opa"),
          variant: "solid"
        }
      ]}
    />
    // <CustomDataTable
    //   data={data}
    //   columns={[
    //     {
    //       key: "project_name",
    //       header: "Project Name",
    //       sortable: true,
    //     },
        // {
        //   key: "is_active",
        //   header: "Status",
        //   sortable: true,
        //   render: (value) => (
        //     <span
        //       className={cn(
        //         "px-2 py-1 rounded-full text-xs font-semibold",
        //         value ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
        //       )}
        //     >
        //       {value ? "Active" : "Inactive"}
        //     </span>
        //   ),
        //   searchableTextExtractor: (value) => (value ? "Active" : "Inactive"),
        // },
        // {
        //   key: "created_at",
        //   header: "Created At",
        //   sortable: true,
        //   render: (value) => formatDate(value as string),
        //   searchableTextExtractor: (value) => formatDate(value as string),
        // },
        // {
        //   key: "updated_at",
        //   header: "Updated At",
        //   sortable: true,
        //   render: (value) => formatDate(value as string),
        //   searchableTextExtractor: (value) => formatDate(value as string),
        // },
    //   ]}
      // tableActions={[
      //   {
      //     label: "Create Scope",
      //     icon: Globe,
      //     onClick: () => console.log("Opa"),
      //     variant: "solid"
      //   }
      // ]}
    // />
  )
}