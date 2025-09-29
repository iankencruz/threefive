// lib/types/auth.ts - Custom auth types for the application
export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  created_at: string;
}

// Auth responses
export interface AuthResponse {
  user: User;
  message: string;
}

// Login request
export interface LoginRequest {
  email: string;
  password: string;
}

// Register request
export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
}

// Password reset request
export interface PasswordResetRequest {
  email: string;
}

// Password reset confirm
export interface PasswordResetConfirmRequest {
  token: string;
  new_password: string;
}

// Change password request
export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

// Validation error from backend
export interface ValidationError {
  field: string;
  message: string;
}

// Error response from backend
export interface ErrorResponse {
  code: string;
  message: string;
  errors?: ValidationError[];
}

// Session info
export interface Session {
  id: string;
  created_at: string;
  updated_at: string;
  expires_at: string;
  ip_address?: string;
  user_agent?: string;
}
