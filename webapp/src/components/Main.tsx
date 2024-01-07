"use client";
import ajaxPromise from "@/lib/ajaxPromise";
import { useQuery } from "@tanstack/react-query";
import React, { useEffect } from "react";
import { DataTable, RespEntity } from "./DataTable";
import { Combobox } from "./Combobox";
import { useRouter, useSearchParams } from "next/navigation";

function Main() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [selectedKind, setSelectedKind] = React.useState<string>("");
  const [sortKey, setSortKey] = React.useState<string>("");
  const [sortDirection, setSortDirection] = React.useState<string>("");

  const { data: kinds } = useQuery({
    queryKey: ["kinds"],
    queryFn: async () => {
      const resp = await ajaxPromise<{ kinds: string[] }>(
        "/api/kinds",
        "GET",
        null
      );
      if (
        searchParams.has("kind") &&
        resp.data.kinds.includes(searchParams.get("kind")!)
      ) {
        setSelectedKind(searchParams.get("kind")!);
      } else {
        setSelectedKind(resp.data.kinds[0]);
      }
      return resp.data.kinds;
    },
  });

  useEffect(() => {
    if (
      searchParams.has("kind") &&
      selectedKind !== searchParams.get("kind") &&
      kinds?.includes(searchParams.get("kind")!)
    ) {
      setSelectedKind(searchParams.get("kind")!);
    }
    if (
      searchParams.has("sortKey") &&
      sortKey !== searchParams.get("sortKey")!
    ) {
      setSortKey(searchParams.get("sortKey")!);
    }
    if (
      searchParams.has("sortDirection") &&
      sortDirection !== searchParams.get("sortDirection")!
    ) {
      setSortDirection(searchParams.get("sortDirection")!);
    }
  }, [searchParams]);
  const { data, isLoading, isError } = useQuery({
    queryKey: ["getEntities", selectedKind, sortKey, sortDirection],
    queryFn: async () => {
      const resp = await ajaxPromise<{ entities: RespEntity[] }>(
        `/api/entities/${selectedKind}/?sortKey=${sortKey}&sortDirection=${sortDirection}`,
        "GET",
        null
      );
      return resp.data.entities;
    },
    enabled: !!selectedKind,
  });
  if (isError) {
    return <div>Error...</div>;
  }
  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <Combobox
        values={kinds || []}
        setValue={(value: string) => {
          setSelectedKind(value);
          setSortKey("");
          setSortDirection("");
          router.push(`/?kind=${value}`);
        }}
        value={selectedKind}
      />
      {data && (
        <DataTable
          data={data}
          onSortChange={(sortings) => {
            setSortKey(sortings[0].id);
            setSortDirection(sortings[0].desc ? "desc" : "asc");
            let path = `/?kind=${selectedKind}&sortKey=${
              sortings[0].id
            }&sortDirection=${sortings[0].desc ? "desc" : "asc"}`;
            router.push(path);
          }}
          defaultSorting={
            sortKey && sortDirection
              ? [
                  {
                    id: sortKey,
                    desc: sortDirection === "desc",
                  },
                ]
              : undefined
          }
        />
      )}
    </div>
  );
}

export default Main;
