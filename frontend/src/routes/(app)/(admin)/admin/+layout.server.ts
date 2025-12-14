// src/routes/admin/+layout.server.ts (New File)
import type { LayoutServerLoad } from "./$types";

// This layout load function runs for: /admin, /admin/media, /admin/pages, etc.
export const load: LayoutServerLoad = async ({ parent }) => {
  // 1. Get the user data from the root layout (+layout.server.ts)
  const { user } = await parent();

  // Pass the user object down to all children pages/layouts
  return { user };
};
