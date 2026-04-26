import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
  // Use the proxy path. Vite will handle the 'http://localhost:8080' part.
  const res = await fetch("/api/admin/dashboard");

  if (res.status === 401) {
    if (typeof window !== "undefined") {
      window.location.href = "http://localhost:8080/login";
    }
    return { dashboard: null };
  }

  if (!res.ok) {
    throw error(res.status, "Backend fetch failed");
  }

  const data = await res.json();
  return { dashboard: data };
};
