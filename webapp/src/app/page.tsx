"use client";
import { DataTable } from "@/components/DataTable";
import { Button } from "@/components/ui/button";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Main from "@/components/Main";
import Image from "next/image";
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
    },
  },
});

export default function Home() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="p-4 space-y-4">
        <Main />
      </div>
    </QueryClientProvider>
  );
}
