// src/routes/+layout.server.ts
import type { LayoutServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: LayoutServerLoad = async ({ fetch }) => {
  try {
    // Try to get current user from backend
    const response = await fetch(`${PUBLIC_API_URL}/auth/me`, {
      credentials: "include",
    });

    if (response.ok) {
      const user = await response.json();
      return { user };
    }
  } catch (error) {
    // Not authenticated or error - that's fine
    console.log("Not authenticated");
  }

  return { user: null };
};
