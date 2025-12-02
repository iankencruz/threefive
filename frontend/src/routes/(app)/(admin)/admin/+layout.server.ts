// src/routes/admin/+layout.server.ts (New File)
import { redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

// This layout load function runs for: /admin, /admin/media, /admin/pages, etc.
export const load: LayoutServerLoad = async ({ parent, url }) => {
  // 1. Get the user data from the root layout (+layout.server.ts)
  const { user } = await parent();

  // 2. CHECK AUTHENTICATION FOR ALL ADMIN ROUTES
  if (!user) {
    // Redirect to login if not authenticated
    redirect(303, `/auth/login`);
  }

  // Pass the user object down to all children pages/layouts
  return { user };
};
