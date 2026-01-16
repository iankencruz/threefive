// frontend/src/routes/(app)/(admin)/admin/contacts/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch, url, parent }) => {
  const { user } = await parent();

  const page = parseInt(url.searchParams.get("page") || "1");
  const limit = 20;
  const status = url.searchParams.get("status");

  try {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });

    // Add status filter if provided
    if (status && status !== "all") {
      params.set("status", status);
    }

    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/admin/contacts?${params.toString()}`,
      {
        credentials: "include",
      },
    );

    if (!response.ok) {
      throw error(response.status, "Failed to fetch contacts");
    }

    const data = await response.json();

    return {
      user,
      contacts: data.contacts || [],
      total: data.total || 0,
      pagination: {
        page: data.offset / data.limit + 1,
        limit: data.limit,
        total_pages: data.total_pages,
        total_count: data.total,
      },
    };
  } catch (err) {
    console.error("Error fetching contacts:", err);
    throw error(500, "Failed to load contacts");
  }
};
