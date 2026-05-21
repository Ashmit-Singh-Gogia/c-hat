import { apiClient } from '../api/client';

export const authService = {
  // Axios call clears the cookie on the backend, then we navigate manually.
  // (Backend handler returns 200 JSON, not a redirect, so window.location alone won't work)
  logout: async () => {
    try { await apiClient.get('/auth/logout'); } catch {}
    window.location.href = '/login';
  }
};