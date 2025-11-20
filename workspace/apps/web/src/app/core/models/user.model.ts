export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  is_active: boolean;
  email_verified: boolean;
  avatar?: string;
  organization?: string;
  created_at: string;
  updated_at: string;
}

export enum UserRole {
  USER = "user",
  ANALYST = "analyst",
  ADMIN = "admin",
}

export interface AuthResponse {
  token: string;
  refresh_token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  organization?: string;
}
