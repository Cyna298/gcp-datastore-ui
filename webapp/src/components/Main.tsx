import ajaxPromise from "@/lib/ajaxPromise";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import { DataTable, RespEntity } from "./DataTable";
import { Combobox } from "./Combobox";

function Main() {
  const [selectedKind, setSelectedKind] = React.useState<string>("");
  const { data: kinds } = useQuery({
    queryKey: ["kinds"],
    queryFn: async () => {
      const resp = await ajaxPromise<{ kinds: string[] }>(
        "/api/kinds",
        "GET",
        null
      );
      setSelectedKind(resp.data.kinds[0]);
      return resp.data.kinds;
    },
  });
  const { data, isLoading, isError } = useQuery({
    queryKey: ["getEntities", selectedKind],
    queryFn: async () => {
      const resp = await ajaxPromise<{ entities: RespEntity[] }>(
        `/api/entities/${selectedKind}/`,
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
        setValue={setSelectedKind}
        value={selectedKind}
      />
      {data && <DataTable data={data} />}
    </div>
  );
}

export default Main;
