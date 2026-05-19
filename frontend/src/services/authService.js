import { apiClient } from '../api/client';

export const authService = {
  getMe: () => apiClient.get('/users/me'),
  // Use relative paths so the Vite proxy forwards these to the backend
  loginUrl: '/api/auth/google',
  logout: () => {
    window.location.href = '/api/auth/google/logout';
  }
};