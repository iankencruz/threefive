import { redirect, type Handle } from "@sveltejs/kit";
import { PUBLIC_API_URL } from "$env/static/public";
import { sequence } from "@sveltejs/kit/hooks";

// Authentication check hook
const authGuard: Handle = async ({ event, resolve }) => {
  const { url, fetch } = event;

  const requiresAuth = url.pathname.startsWith("/admin");

  if (requiresAuth) {
    const path = url.pathname;

    console.log(`\n 🔐 [AUTH CHECK]  Path: ${path} \n`);

    try {
      const response = await fetch(`${PUBLIC_API_URL}/auth/me`, {
        credentials: "include",
        headers: {
          cookie: event.request.headers.get("cookie") ?? "",
        },
      });

      if (!response.ok) {
        throw redirect(303, "/auth/login");
      }

      const user = await response.json();
      event.locals.user = user;

      console.log(
        `\n ✅ [AUTH OK]      User: ${user.email ?? "unknown"}  • Path: ${path} \n`,
      );
    } catch (err) {
      // Preserve redirects
      if (err instanceof Response && err.status === 303) {
        throw err;
      }

      // Narrow unknown error safely
      const message =
        err instanceof Error
          ? err.message
          : typeof err === "string"
            ? err
            : JSON.stringify(err);

      console.error(
        `❌ [AUTH ERROR] • Not authenticated. Redirecting → /auth/login \n • Path: ${path}   • Error: ${message}`,
      );

      throw redirect(303, "/auth/login");
    }
  }

  return resolve(event);
};

// Export the handle function
export const handle = sequence(authGuard);
