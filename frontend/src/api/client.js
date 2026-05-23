import axios from 'axios';

// Use the Vite dev server proxy by requesting the frontend-origin `/api` path.
// This keeps requests same-origin so HttpOnly cookies are handled by the browser
// without requiring cross-origin credential handling in development.
export const apiClient = axios.create({
  baseURL: '/api',
  withCredentials: true, // Crucial for HTTP-Only Cookies
});