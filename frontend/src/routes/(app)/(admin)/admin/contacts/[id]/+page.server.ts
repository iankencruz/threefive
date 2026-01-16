// frontend/src/routes/(app)/(admin)/admin/contacts/[id]/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad, Actions } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch, params, parent }) => {
  const { user } = await parent();

  try {
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/admin/contacts/${params.id}`,
      {
        credentials: "include",
      },
    );

    if (!response.ok) {
      if (response.status === 404) {
        throw error(404, "Contact not found");
      }
      throw error(response.status, "Failed to fetch contact");
    }

    const contact = await response.json();

    return {
      user,
      contact,
    };
  } catch (err) {
    console.error("Error fetching contact:", err);
    if (
      typeof err === "object" &&
      err !== null &&
      "status" in err &&
      typeof err.status === "number"
    ) {
      throw err;
    }
    throw error(500, "Failed to load contact");
  }
};

export const actions: Actions = {
  updateStatus: async ({ request, params, fetch }) => {
    const formData = await request.formData();
    const status = formData.get("status");

    try {
      const response = await fetch(
        `${PUBLIC_API_URL}/api/v1/admin/contacts/${params.id}/status`,
        {
          method: "PATCH",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ status }),
        },
      );

      if (!response.ok) {
        return {
          success: false,
          error: "Failed to update status",
        };
      }

      return {
        success: true,
      };
    } catch (err) {
      console.error("Error updating status:", err);
      return {
        success: false,
        error: "Failed to update status",
      };
    }
  },

  delete: async ({ params, fetch }) => {
    try {
      const response = await fetch(
        `${PUBLIC_API_URL}/api/v1/admin/contacts/${params.id}`,
        {
          method: "DELETE",
          credentials: "include",
        },
      );

      if (!response.ok) {
        return {
          success: false,
          error: "Failed to delete contact",
        };
      }

      return {
        success: true,
        deleted: true,
      };
    } catch (err) {
      console.error("Error deleting contact:", err);
      return {
        success: false,
        error: "Failed to delete contact",
      };
    }
  },
};
