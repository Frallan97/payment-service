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

  static logout(): void {
    this.clearToken();
    window.location.href = '/login';
  }

  static getAuthHeader() {
    const token = this.getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }
}
