import { FileText, BarChart3 } from "lucide-react";
import { ROUTES } from "@/types/routes";

export const navItems = [
  {
    title: "System Analysis",
    icon: BarChart3,
    url: ROUTES.SYSTEM_ANALYSIS,
  },
  {
    title: "Documentation",
    icon: FileText,
    url: ROUTES.DOCUMENTATION,
  },
];

export const secondaryNavItems = [];
