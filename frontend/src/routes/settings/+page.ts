import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
  // Use the proxy path. Vite will handle the 'http://localhost:8080' part.
  const res = await fetch("/api/admin/system-config");

  if (res.status === 401) {
    throw redirect(302, "http://localhost:8080/login");
  }

  if (!res.ok) {
    throw error(res.status, "Backend fetch failed");
  }

  const data = await res.json();
  return { config: data };
};
