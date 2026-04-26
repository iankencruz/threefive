// src/routes/+layout.ts
export const ssr = false;
export const prerender = false; // Usually false for dynamic apps, or true if pages are static

import type { LayoutLoad } from "./$types";

// This runs ONLY ONCE when the app starts up in the browser
const configPromise = fetch("/api/admin/system-config/application_name")
  .then((res) => res.json())
  .catch(() => ({ application_name: "Default App" }));

export const load: LayoutLoad = async () => {
  const config = await configPromise;
  return {
    appName: config.value,
  };
};
