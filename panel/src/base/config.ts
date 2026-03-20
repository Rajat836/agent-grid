export const Config = {
  appName: process.env.NEXT_PUBLIC_APP_NAME || "MyApp",
  apiHost:
    process.env.NEXT_PUBLIC_API_HOST || process.env.NEXT_PUBLIC_BASE_URL || "",
};
