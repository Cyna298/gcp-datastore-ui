import * as React from "react";
import { CaretSortIcon, ChevronDownIcon } from "@radix-ui/react-icons";
import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  useReactTable,
} from "@tanstack/react-table";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { Badge } from "./ui/badge";

export type RespEntity = {
  [key: string]: {
    name: string;
    value: string;
    type: string;
    indexed: boolean;
  };
};
type BadgeStyled = {
  label: string;
  className: string;
};
const typesToClassNames: Record<string, BadgeStyled> = {
  string: {
    label: "S",
    className: "bg-blue-500 text-white",
  },
  int64: {
    label: "N",
    className: "bg-green-500 text-white",
  },
  bool: {
    label: "B",
    className: "bg-yellow-500 text-black",
  },
  "<nil>": {
    label: "N",
    className: "bg-red-500 text-white",
  },

  "time.Time": {
    label: "T",
    className: "bg-indigo-500 text-white",
  },
  "[]interface {}": {
    label: "A",
    className: "bg-indigo-500 text-white",
  },
};

function getColumns(entities: RespEntity[]): ColumnDef<RespEntity>[] {
  const columns: ColumnDef<RespEntity>[] = [];
  const keysSet = new Set<string>();
  for (const entity of entities) {
    for (const key of Object.keys(entity)) {
      keysSet.add(key);
    }
  }
  const keys = Array.from(keysSet).sort((a, b) => {
    if (a == "key") {
      return -1;
    }
    if (b == "key") {
      return 1;
    }

    return a.localeCompare(b);
  });
  keys.forEach((key) => {
    const column: ColumnDef<RespEntity> = {
      accessorKey: key,
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            {key}
            {!column.getIsSorted() ? (
              <CaretSortIcon className="ml-2 h-4 w-4" />
            ) : (
              <ChevronDownIcon
                className={`ml-2 h-4 w-4 ${
                  column.getIsSorted() === "asc" ? "transform rotate-180" : ""
                }`}
              />
            )}
          </Button>
        );
      },

      cell: ({ row }) => {
        const value = row.original[key]?.value;
        const type = value ? row.original[key]?.type : null;
        const badge = type ? typesToClassNames?.[type] : null;
        return (
          <HoverCard>
            <HoverCardTrigger>
              <div className="flex items-center space-x-2">
                <div className="truncate max-w-96">{value || "-"}</div>
                {badge ? (
                  <Badge className={`text-xs flex-none ${badge.className}`}>
                    {badge.label}
                  </Badge>
                ) : (
                  type && <Badge className="text-xs flex-none">{type}</Badge>
                )}
              </div>
            </HoverCardTrigger>
            <HoverCardContent className="w-full" align="start">
              <div className="flex flex-wrap">{value}</div>
            </HoverCardContent>
          </HoverCard>
        );
      },
    };
    columns.push(column);
  });
  return columns;
}
export function DataTable({
  data,
  onSortChange,
  defaultSorting,
}: {
  data: RespEntity[];
  onSortChange: (sortingState: SortingState) => void;
  defaultSorting?: SortingState;
}) {
  const [sorting, setSorting] = React.useState<SortingState>(
    defaultSorting ?? []
  );
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState({});

  const columns = data ? getColumns(data) : [];

  const table = useReactTable({
    data: data ?? [],
    columns: columns,
    onSortingChange: (sortFn) => {
      //check if sorting is callable
      if (typeof sortFn === "function") {
        onSortChange(sortFn(sorting));
      }

      setSorting(sortFn);
    },
    manualSorting: true,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  });

  return (
    <div className="w-full">
      <div className="flex items-center py-4">
        {/* <Input
          placeholder="Filter emails..."
          value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("email")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        /> */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="ml-auto">
              Columns <ChevronDownIcon className="ml-2 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter((column) => column.getCanHide())
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                );
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <div className="flex-1 text-sm text-muted-foreground">
          {table.getFilteredSelectedRowModel().rows.length} of{" "}
          {table.getFilteredRowModel().rows.length} row(s) selected.
        </div>
        <div className="space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
