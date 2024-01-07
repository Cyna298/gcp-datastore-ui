"use client";
import ajaxPromise from "@/lib/ajaxPromise";
import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import React, { useEffect } from "react";
import { DataTable, RespEntity } from "./DataTable";
import { Combobox } from "./Combobox";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "./ui/button";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ReloadIcon,
} from "@radix-ui/react-icons";

function Main() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [selectedKind, setSelectedKind] = React.useState<string>("");
  const [sortKey, setSortKey] = React.useState<string>("");
  const [sortDirection, setSortDirection] = React.useState<string>("");
  const [selectedPage, setSelectedPage] = React.useState<number>(0);

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
  const fetchEntities = async ({ pageParam = "" }) => {
    const resp = await ajaxPromise<{
      entities: RespEntity[];
      nextCursor: string;
    }>(
      `/api/entities/${selectedKind}/?sortKey=${sortKey}&sortDirection=${sortDirection}&cursor=${pageParam}&limit=50`,
      "GET",
      null
    );
    return resp.data;
  };
  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isLoading,
    isError,
    isFetchingNextPage,
    status,
  } = useInfiniteQuery({
    queryKey: ["getEntities", selectedKind, sortKey, sortDirection],
    queryFn: fetchEntities,
    getNextPageParam: (lastPage) => {
      if (lastPage.nextCursor) {
        return lastPage.nextCursor;
      }
    },
    enabled: !!selectedKind,
    initialPageParam: "",
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
          setSelectedPage(0);
          router.push(`/?kind=${value}`);
        }}
        value={selectedKind}
      />
      {data && (
        <div className="mt-2">
          <DataTable
            data={data.pages[selectedPage].entities}
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
          <div className="w-full flex space-x-2 items-center justify-end mt-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setSelectedPage((prev) => prev - 1)}
              disabled={isFetchingNextPage || selectedPage === 0}
            >
              <ChevronLeftIcon className="h-4 w-4" />
              <span className="ml-2">Previous</span>
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                if (selectedPage === data.pages.length - 1) {
                  fetchNextPage().then((res) => {
                    setSelectedPage((prev) => prev + 1);
                  });
                } else {
                  setSelectedPage((prev) => prev + 1);
                }
              }}
              disabled={
                selectedPage === data.pages.length - 1
                  ? !hasNextPage || isFetchingNextPage
                  : false
              }
            >
              <span className="mr-2">Next</span>
              {isFetchingNextPage ? (
                <ReloadIcon className="h-4 w-4 animate-spin" />
              ) : (
                <ChevronRightIcon className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Main;
