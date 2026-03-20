"use client";

import Image from "next/image";
import { usePathname } from "next/navigation";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { navItems } from "@/lib/sidebar";

export function AppSidebar() {
  const pathname = usePathname();

  return (
    <Sidebar collapsible="none" className="border-r border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-slate-900">
      <SidebarHeader className="border-b border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-slate-900 py-3">
        <div className="flex items-center justify-between p-2">
          <span className="text-lg font-bold whitespace-nowrap flex items-center gap-2 text-gray-900 dark:text-gray-100">
            <Image src="/map.png" alt="Logo" width={24} height={24} />
            <span>Ontology</span>
          </span>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    asChild
                    isActive={
                      pathname === item.url ||
                      pathname.startsWith(`${item.url}/`)
                    }
                    tooltip={item.title}
                    className="transition-all hover:bg-gray-100 dark:hover:bg-gray-800"
                  >
                    <a href={item.url} className="group">
                      <item.icon className="h-4 w-4 group-hover:text-gray-700 dark:group-hover:text-gray-300 transition-colors" />
                      <span className="group-hover:text-gray-700 dark:group-hover:text-gray-300 transition-colors">{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {/* No secondary nav items */}
      </SidebarContent>

      <SidebarFooter className="border-t border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-slate-900">
        {/* Footer content can be added here */}
      </SidebarFooter>
    </Sidebar>
  );
}
