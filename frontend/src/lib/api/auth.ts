import { PUBLIC_API_URL } from "$env/static/public";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  PasswordResetRequest,
  PasswordResetConfirmRequest,
  ChangePasswordRequest,
  User,
  Session,
  ErrorResponse,
} from "$types/auth";

// Base fetch wrapper with error handling
async function fetchAPI<T>(url: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(url, {
    credentials: "include", // Important for cookies
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  const data = await response.json();

  if (!response.ok) {
    throw data as ErrorResponse;
  }

  return data as T;
}

// Auth API functions
export const authApi = {
  // Register new user
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return fetchAPI<AuthResponse>(`${PUBLIC_API_URL}/auth/register`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  // Login user
  async login(data: LoginRequest): Promise<AuthResponse> {
    return fetchAPI<AuthResponse>(`${PUBLIC_API_URL}/auth/login`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  // Logout current session
  async logout(): Promise<{ message: string }> {
    return fetchAPI(`${PUBLIC_API_URL}/auth/logout`, {
      method: "POST",
    });
  },

  // Logout all sessions
  async logoutAll(): Promise<{ message: string }> {
    return fetchAPI(`${PUBLIC_API_URL}/auth/logout-all`, {
      method: "POST",
    });
  },

  // Get current user
  async me(): Promise<User> {
    return fetchAPI<User>(`${PUBLIC_API_URL}/auth/me`);
  },

  // Change password
  async changePassword(
    data: ChangePasswordRequest,
  ): Promise<{ message: string }> {
    return fetchAPI("/auth/change-password", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  // Request password reset
  async requestPasswordReset(
    data: PasswordResetRequest,
  ): Promise<{ message: string }> {
    return fetchAPI("/auth/request-password-reset", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  // Reset password with token
  async resetPassword(
    data: PasswordResetConfirmRequest,
  ): Promise<{ message: string }> {
    return fetchAPI("/auth/reset-password", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  // Get user sessions
  async getSessions(): Promise<{ sessions: Session[] }> {
    return fetchAPI("/auth/sessions");
  },

  // Revoke specific session
  async revokeSession(sessionId: string): Promise<{ message: string }> {
    return fetchAPI(`/auth/sessions/${sessionId}`, {
      method: "DELETE",
    });
  },

  // Refresh session
  async refreshSession(): Promise<{ message: string }> {
    return fetchAPI("/auth/refresh", {
      method: "POST",
    });
  },
};
