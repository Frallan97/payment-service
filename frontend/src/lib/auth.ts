export class AuthService {
  private static readonly TOKEN_KEY = 'auth_token';

  static isAuthenticated(): boolean {
    return !!localStorage.getItem(this.TOKEN_KEY);
  }

  static getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  static setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  static clearToken(): void {
    localStorage.removeItem(this.TOKEN_KEY);
  }

  static async login(email: string, password: string): Promise<boolean> {
    // Simple mock login for now
    // TODO: Replace with actual API call
    if (email && password) {
      this.setToken('mock-token');
      return true;
    }
    return false;
  }

  static logout(): void {
    this.clearToken();
    window.location.href = '/login';
  }

  static getAuthHeader() {
    const token = this.getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }
}
