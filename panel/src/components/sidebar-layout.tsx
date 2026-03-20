"use client";

import { AppSidebar } from "@/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";

export function SidebarLayout({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <SidebarProvider defaultOpen={true} className="h-svh overflow-hidden" style={{ "--sidebar-width": "16rem" } as React.CSSProperties}>
      <AppSidebar />
      <SidebarInset className={`overflow-hidden flex flex-col ${className || ""}`}>
        <div className="h-full w-full overflow-y-auto">{children}</div>
      </SidebarInset>
    </SidebarProvider>
  );
}
