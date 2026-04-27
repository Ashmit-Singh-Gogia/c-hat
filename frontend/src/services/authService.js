import { apiClient } from '../api/client';

export const authService = {
  getMe: () => apiClient.get('/users/me'),
  loginUrl: 'http://localhost:8082/api/auth/google',
  logout: () => {
    document.cookie = "jwt_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    window.location.href = '/login';
  }
};