const AUTH_SERVICE_URL = import.meta.env.VITE_AUTH_SERVICE_URL || "https://auth.vibeoholic.com";
const PAYMENT_SERVICE_URL = import.meta.env.VITE_PAYMENT_SERVICE_URL || "https://payments.vibeoholic.com";

export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  is_super_admin?: boolean;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

export class AuthService {
  private static ACCESS_TOKEN_KEY = "access_token";
  private static REFRESH_TOKEN_KEY = "refresh_token";
  private static USER_KEY = "user";

  // OAuth Login - redirect to auth-service
  static initiateOAuthLogin(): void {
    const redirectUri = window.location.origin;
    const authUrl = `${AUTH_SERVICE_URL}/api/auth/google/login?redirect_uri=${encodeURIComponent(redirectUri)}`;
    window.location.href = authUrl;
  }

  // Handle OAuth callback
  static async handleOAuthCallback(accessToken: string): Promise<User> {
    // Store access token
    // Note: refresh_token is stored as HTTP-only cookie by auth service
    localStorage.setItem(this.ACCESS_TOKEN_KEY, accessToken);

    // Fetch user info from auth-service
    const user = await this.fetchCurrentUser();
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));

    return user;
  }

  // Fetch current user from auth-service
  private static async fetchCurrentUser(): Promise<User> {
    const token = this.getAccessToken();
    if (!token) {
      throw new Error("No access token available");
    }

    const response = await fetch(`${AUTH_SERVICE_URL}/api/auth/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      throw new Error("Failed to fetch user info");
    }

    return response.json();
  }

  static logout(): void {
    localStorage.removeItem(this.ACCESS_TOKEN_KEY);
    localStorage.removeItem(this.REFRESH_TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
    window.location.href = "/login";
  }

  static getAccessToken(): string | null {
    return localStorage.getItem(this.ACCESS_TOKEN_KEY);
  }

  static getRefreshToken(): string | null {
    return localStorage.getItem(this.REFRESH_TOKEN_KEY);
  }

  static getUser(): User | null {
    const userStr = localStorage.getItem(this.USER_KEY);
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }

  static isAuthenticated(): boolean {
    const token = this.getAccessToken();
    const user = this.getUser();
    return !!(token && user);
  }

  static async fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
    const token = this.getAccessToken();

    if (!token) {
      throw new Error("Not authenticated");
    }

    const headers = new Headers(options.headers);
    headers.set("Authorization", `Bearer ${token}`);

    return fetch(url, {
      ...options,
      headers,
    });
  }

  // Refresh access token using refresh token (from HTTP-only cookie)
  static async refreshAccessToken(): Promise<string> {
    const response = await fetch(`${AUTH_SERVICE_URL}/api/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include", // Send HTTP-only cookie
    });

    if (!response.ok) {
      // Refresh failed, clear tokens and redirect to login
      this.logout();
      throw new Error("Token refresh failed");
    }

    const data = await response.json();
    localStorage.setItem(this.ACCESS_TOKEN_KEY, data.access_token);

    return data.access_token;
  }
}

export { PAYMENT_SERVICE_URL };
