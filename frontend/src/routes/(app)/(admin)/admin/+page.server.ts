// src/routes/admin/+page.server.ts

// This file is empty because all it needs to do is exist.
// Its existence makes /admin a valid route, allowing the redirect
// in +layout.server.ts to execute cleanly.
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ parent, url }) => {
  const { user } = await parent();
  // Redirect the root /admin path to /admin/dashboard
  // Check if the current pathname is EXACTLY /admin (or /admin/)
  if (url.pathname === "/admin" || url.pathname === "/admin/") {
    redirect(303, `/admin/dashboard`);
  }
};
