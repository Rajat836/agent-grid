
export const ROUTES = {
  HOME: "/",
  SYSTEM_ANALYSIS: "/system-analysis",
  DOCUMENTATION: "/documentation",
} as const;

// Type for all valid routes
export type Route = (typeof ROUTES)[keyof typeof ROUTES] | string;
