import { apiClient } from '../api/client';

export const authService = {
  getMe: () => apiClient.get('/users/me'),
  loginUrl: 'http://localhost:8082/api/auth/google',
  logout: () => {
    window.location.href = 'http://localhost:8082/api/auth/google/logout';
  }
};