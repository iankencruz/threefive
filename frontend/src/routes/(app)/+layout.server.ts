import type { LayoutServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: LayoutServerLoad = async ({ locals, fetch }) => {
  // For admin routes, user is already set in locals by hooks
  if (locals.user) {
    return { user: locals.user };
  }

  // For non-admin routes, still try to fetch user (optional auth)
  try {
    const response = await fetch(`${PUBLIC_API_URL}/auth/me`, {
      credentials: "include",
    });

    if (response.ok) {
      const user = await response.json();
      return { user };
    }
  } catch (error) {
    console.log("Not authenticated");
  }

  return { user: null };
};
